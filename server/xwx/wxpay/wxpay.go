package wxpay

import (
	"context"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xdatabase/xredis"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"time"
)

var Cli *wechat.ClientV3

var Cfg *Config

type Config struct {
	AppID      string `yaml:"appId"`      // appid 小程序ID
	MchId      string `yaml:"mchId"`      // 商户ID 或者服务商模式的 sp_mchid
	MchKey     string `yaml:"mchKey"`     // appKey apiV3Key
	PrivateKey string `yaml:"privateKey"` // 私钥 apiclient_key.pem 读取后的内容
	NotifyUrl  string `yaml:"notifyUrl"`  // 回调地址
	SerialNo   string `yaml:"serialNo"`   //  商户证书的证书序列号
	MchName    string `yaml:"mchName"`    // 商户名称 支付备注使用默认
}

func NotifyUrl() string {
	return Cfg.NotifyUrl
}

func InitWith() *wechat.ClientV3 {
	config, err := x.Config[Config]("wx.pay")
	if err != nil {
		xlog.Panicf(context.Background(), "wx.pay config error: %v", err)
	}
	if config.MchName == "" {
		config.MchName = "商品"
	}
	return Init(config)
}

func Init(config *Config) *wechat.ClientV3 {
	Cfg = config
	if Cfg.AppID == "" {
		xlog.Errorf(nil, "wx pay xconfig is empty init error")
		return nil
	}
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容
	client, err := wechat.NewClientV3(Cfg.MchId, Cfg.SerialNo, Cfg.MchKey, Cfg.PrivateKey)
	if err != nil {
		xlog.Warnf(nil, "wx pay init error %s", err)
		panic(err)
	}
	err = client.AutoVerifySign() // 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	if err != nil {
		xlog.Errorf(nil, "wx pay AutoVerifySign error %s", err)
		panic(err)
	}
	// client.DebugSwitch = gopay.DebugOn // 调试开启
	xlog.Infof(nil, "wx pay init success %s==>%s", Cfg.AppID, Cfg.MchId)
	Cli = client
	return client
}

type PlaceOrder struct {
	Amount int32     // 金额分
	Openid string    // openId
	PaySn  int64     // 支付单号
	Desc   string    // 商品描述
	Expire time.Time //
}

// TransactionJsapiPrepayId 获取预支付交易会话标识。用于后续接口调用中使用，该值有效期为2小时
func TransactionJsapiPrepayId(ctx context.Context, place *PlaceOrder) (string, error) {

	if place == nil {
		return "", errors.New("pay params error")
	}
	// 缓存90分钟小时
	client, err := xredis.Get()
	if err != nil {
		xlog.Errorf(ctx, "wx pay get xcache error %s", err)
		return "", err
	}
	cacheKey := fmt.Sprintf("wx_pay_prepayid_%d", place.PaySn)
	get := client.Get(ctx, cacheKey)
	result, err := get.Result()
	if err != nil && !errors.Is(err, xredis.Nil) {
		xlog.Errorf(ctx, "wx pay get xcache prepayid error %s", err)
		return "", err
	}
	if result != "" {
		return result, nil
	}

	expire := place.Expire.Format(time.RFC3339)

	if place.Amount <= 0 {
		return "", errors.New("amount error")
	}
	if place.Openid == "" {
		return "", errors.New("openid error")
	}
	if place.PaySn <= 0 {
		return "", errors.New("pay sn error")
	}
	if place.Desc == "" {
		place.Desc = Cfg.MchName
	}
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("appid", Cfg.AppID).
		Set("mchid", Cfg.MchId).
		Set("description", place.Desc).                      // 商品描述
		Set("out_trade_no", fmt.Sprintf("%d", place.PaySn)). // 商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一 6-32位
		// 订单失效时间，遵循rfc3339标准格式，格式为yyyy-MM-DDTHH:mm:ss+TIMEZONE，yyyy-MM-DD表示年月日，T出现在字符串中，表示time元素的开头，HH:mm:ss表示时分秒，TIMEZONE表示时区（+08:00表示东八区时间，领先UTC8小时，即北京时间）。例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日 13点29分35秒。
		Set("time_expire", expire).
		// 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。 公网域名必须为https，如果是走专线接入，使用专线NAT IP或者私有回调域名可使用http
		Set("notify_url", Cfg.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", place.Amount). // 订单总金额，单位为分。
							Set("currency", "CNY") // CNY：人民币，境内商户号仅支持人民币。
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", place.Openid)
		})

	wxRsp, err := Cli.V3TransactionJsapi(ctx, bm)
	if err != nil {
		xlog.Errorf(ctx, "wx pay error %v", err)
		return "", err
	}
	xlog.Debugf(ctx, "build %#v wxRsp: %#v", bm, wxRsp.Response)
	// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_1.shtml
	if wxRsp.Code == wechat.Success {
		xlog.Debugf(ctx, "wxRsp: %#v", wxRsp.Response)
		id := wxRsp.Response.PrepayId
		client.Set(ctx, cacheKey, id, 90*time.Minute)
		return id, nil // 预支付交易会话标识。用于后续接口调用中使用，该值有效期为2小时
	} else {
		xlog.Errorf(ctx, "build %#v wxRsp: %#v", bm, wxRsp.Response)
	}

	return "", errors.New("wx pay error")
}

func GetMiniPaySign(prepayid string) (*wechat.AppletParams, error) {
	// 小程序
	return Cli.PaySignOfApplet(Cfg.AppID, prepayid)
}
