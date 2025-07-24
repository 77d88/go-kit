package aliyunaddress

import (
	"encoding/json"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xcache"
	"github.com/77d88/go-kit/plugins/xlog"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var client *sdk.Client

func Init(ak, sk string) xhs.Handler {
	credentialsProvider := credentials.NewStaticAKCredentialsProvider(ak, sk)
	c, err := sdk.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), credentialsProvider)
	if err != nil {
		xlog.Fatalf(nil, "aliyun address client init failed: %s", err)
	}
	client = c
	return func(ctx *xhs.Ctx) (interface{}, error) {
		return Run(ctx)
	}
}

// handler 地址标准化 /address/standardizeAddress
func handler(c *xhs.Ctx, r request) (interface{}, error) {

	var res response
	err := xcache.Once("aliyun_address_cache_"+r.Text, &res, time.Minute*30, func() (interface{}, error) {
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
			c.Fatalf(err, "解析失败")
		}
		var obj httpRes
		c.Fatalf(json.Unmarshal([]byte(res.GetHttpContentString()), &obj), "解析失败")
		var data httpInfo
		c.Fatalf(json.Unmarshal([]byte(obj.Data), &data), "解析失败")
		c.Fatalf(data.Status != "OK", "解析失败")

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

	return res, err

}

// Run 地址标准化
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	c.ShouldBind(&r)
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
