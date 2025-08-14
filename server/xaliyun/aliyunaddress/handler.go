package aliyunaddress

import (
	"context"
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xcache"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var client *sdk.Client
var once sync.Once

func Init() *sdk.Client {
	once.Do(func() {
		xlog.Infof(context.Background(), "aliyunaddress.Init")

		ak := x.ConfigString("aliyun.address.ak")
		sk := x.ConfigString("aliyun.address.ak")
		if ak == "" || sk == "" {
			xlog.Fatalf(nil, "aliyun address ak or sk is empty")
		}
		credentialsProvider := credentials.NewStaticAKCredentialsProvider(ak, sk)
		c, err := sdk.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), credentialsProvider)
		if err != nil {
			xlog.Fatalf(nil, "aliyun address client init failed: %s", err)
		}
		client = c
		if client != nil {
			x.Use(func() *sdk.Client { return client }, true)
		}
	})

	return client
}

// handler 地址标准化 /address/standardizeAddress
func handler(c *xhs.Ctx, r request) (interface{}, error) {

	return xcache.Once("aliyun_address_cache_"+r.Text, time.Minute*30, func() (interface{}, error) {
		// 构造一个公共请求
		request := requests.NewCommonRequest()
		// 设置请求方式
		request.Method = "POST"
		// 指定产品
		request.Product = "address-purification"
		// 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
		request.Domain = "address-purification.cn-hangzhou.aliyuncs.com"
		// 指定产品版本
		request.Version = "2019-11-18"
		// 指定接口名
		request.ApiName = "ExtractExpress"
		// 设置参数值
		request.QueryParams["ServiceCode"] = "addrp"
		// 设置参数值
		request.QueryParams["AppKey"] = "d6hvepuyxu4b"
		// 设置参数值
		request.QueryParams["Text"] = r.Text
		// 设置参数值
		request.QueryParams["DefaultProvince"] = "浙江省"
		// 设置参数值
		request.QueryParams["DefaultCity"] = "杭州市"
		// 设置参数值
		request.QueryParams["DefaultDistrict"] = "余杭区"
		// 把公共请求转化为acs请求
		request.TransToAcsRequest()

		res, err := client.ProcessCommonRequest(request)
		if err != nil {
			xlog.Errorf(nil, "aliyun address client request failed: %s", err)
			return nil, xerror.New("处理失败")
		}
		obj, err := xparse.FromJSONNew[httpRes](res.GetHttpContentString())
		if err != nil {
			xlog.Errorf(nil, "aliyun address client request failed: %s", err)
			return nil, xerror.New("解析响应失败")
		}
		data, err := xparse.FromJSONNew[httpInfo](obj.Data)
		if err != nil {
			xlog.Errorf(nil, "aliyun address client request failed: %s", err)
			return nil, xerror.New("解析数据失败")
		}
		if data.Status != "OK" {
			xlog.Errorf(nil, "aliyun address client request failed: %s", err)
			return nil, xerror.New("解析状态错误")
		}

		return &response{
			UserName:     data.ExpressExtract.Per,
			ProvinceName: data.ExpressExtract.Prov,
			CityName:     data.ExpressExtract.City,
			CountyName:   data.ExpressExtract.District,
			StreetName:   data.ExpressExtract.PoiInfo,
			DetailInfo:   data.ExpressExtract.AddrInfo,
			TelNumber:    data.ExpressExtract.Tel,
		}, nil
	})
}

// Run 地址标准化
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	err := c.ShouldBind(&r)
	if err != nil {
		return nil, err
	}
	return handler(c, r)
}

type request struct {
	Text string `json:"text" form:"text"`
}

type response struct {
	UserName     string `json:"userName"`
	TelNumber    string `json:"telNumber"`
	ProvinceName string `json:"provinceName"`
	CityName     string `json:"cityName"`
	CountyName   string `json:"countyName"`
	DetailInfo   string `json:"detailInfo"`
	StreetName   string `json:"streetName"`
}

type httpRes struct {
	RequestId string `json:"RequestId"`
	Data      string `json:"Body"`
}
type httpInfo struct {
	ExpressExtract struct {
		HouseInfo string `json:"house_info"`
		PoiInfo   string `json:"poi_info"`
		Town      string `json:"town"`
		City      string `json:"city"`
		District  string `json:"district"`
		Tel       string `json:"tel"`
		AddrInfo  string `json:"addr_info"`
		Per       string `json:"per"`
		Prov      string `json:"prov"`
	} `json:"express_extract"`
	Status   string `json:"status"`
	TimeUsed struct {
		Rt struct {
			BasicChunking  float64 `json:"basic_chunking"`
			Segment        float64 `json:"segment"`
			AddressCorrect float64 `json:"address_correct"`
			Complete       float64 `json:"complete"`
			ExpressExtract float64 `json:"express_extract"`
			AddressSearch  float64 `json:"address_search"`
			Structure      float64 `json:"structure"`
		} `json:"rt"`
		Start float64 `json:"start"`
	} `json:"time_used"`
}
