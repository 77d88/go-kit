package aliyunaddress

//import (
//	"encoding/json"
//	"github.com/77d88/go-kit/plugins/xapi"
//	"github.com/77d88/go-kit/plugins/xdb"
//	"testing"
//)
//
//func Test_Run(t *testing.T) {
//	xapi.InitTestConfig()
//	xdb.Init()
//	str := `{"RequestId":"96068AA0-5C55-5AE1-9F8A-83F33FE665E0","Body":"{\"express_extract\":{\"house_info\":\"1号楼2单元1102室\",\"poi_info\":\"北三环西路19号院\",\"town\":\"瓶窑镇\",\"city\":\"杭州市\",\"district\":\"余杭区\",\"tel\":\"15380991239\",\"addr_info\":\"北三环西路19号院1号楼2单元1102室\",\"per\":\"张三\",\"prov\":\"浙江省\"},\"status\":\"OK\",\"time_used\":{\"rt\":{\"basic_chunking\":0.03929591178894043,\"segment\":0.0063245296478271484,\"address_correct\":0.003167390823364258,\"complete\":0.0002808570861816406,\"express_extract\":0.00003528594970703125,\"address_search\":0.38698649406433105,\"structure\":0.00013375282287597656},\"start\":1732256605.8475428}}"}`
//
//	var obj httpRes
//	err := json.Unmarshal([]byte(str), &obj)
//	var data httpInfo
//	err = json.Unmarshal([]byte(obj.Data), &data)
//
//	t.Log(err)
//
//	res := response{
//		UserName:     data.ExpressExtract.Per,
//		ProvinceName: data.ExpressExtract.Prov,
//		CityName:     data.ExpressExtract.City,
//		CountyName:   data.ExpressExtract.District,
//		StreetName:   data.ExpressExtract.PoiInfo,
//		DetailInfo:   data.ExpressExtract.AddrInfo,
//		TelNumber:    data.ExpressExtract.Tel,
//	}
//	t.Log(res)
//
//	// Run(xtest.ApiContext(request{
//	// 	Text: "张三 15380991239 北京市朝阳区北三环西路19号院1号楼2单元1102室",
//	// }))
//}
