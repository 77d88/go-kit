package xdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/xlog"
	"gorm.io/gorm"
)

type MuDbUser struct {
	BaseModel
	Username string
	WxOpenId string
}

func (m *MuDbUser) TableName() string {
	return "s_user"
}

type MuDbProduct struct {
	BaseModel
}

func (m *MuDbProduct) Limit() (offset, limit int) {
	return 1, 10
}

func (m *MuDbProduct) TableName() string {
	return "s_product"
}

func idMaxScope(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id > ?", id)
	}
}

func TestBaseFunc(t *testing.T) {
	New(&Config{
		Dns:    FastDsn("127.0.0.1", 5432, "postgres", "jerry123!", "zyv2"),
		Logger: true,
	})
	//var dbusers []MuDbUser

	db, _ := GetDB()
	//take, err := NewParams[MuDbUser]().
	//	Eq("id", 1).
	//	ILike("username", "超").
	//	BuildC(db).Take(context.Background())
	//if err != nil {
	//	t.Error(err)
	//}
	//xlog.Infof(nil, "take %+v", take)
	// 获取底层 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// 预热
	sqlDB.Ping()

	// 测试
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := sqlDB.Exec("SELECT 1")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Attempt %d: %v ms\n", i+1, time.Since(start).Milliseconds())
	}

	page := FindPage[MuDbUser](db, &MuDbProduct{}, true)
	//find, err := gorm.G[MuDbUser](where).Find(context.Background())
	xlog.Infof(nil, "take %+v", page)
	ix := []int32{1, 2}
	NewInt8Array(ix...)

}

func TestMuDb(t *testing.T) {
	New(&Config{
		Dns: FastDsn("127.0.0.1", 5432, "postgres", "jerry123!", "zyv2"),
	})
	New(&Config{
		Dns:        FastDsn("127.0.0.1", 5432, "postgres", "jerry123!", "gamev2"),
		DbLinkName: "game",
	})

	var db1Users []MuDbUser
	err := C(context.Background()).Raw(`select id,username  "name" from s_user order by id desc limit 3`).Find(&db1Users)
	if err != nil {
		t.Error(err)
	}

	xarray.ForEach(db1Users, func(index int, item MuDbUser) {
		fmt.Printf("%+v\n", item)
	})

	var db2Users []MuDbUser
	err = C(context.Background(), "game").Raw("select * from g_user order by id desc limit 3").Find(&db2Users)
	if err != nil {
		t.Error(err)
	}
	xarray.ForEach(db2Users, func(index int, item MuDbUser) {
		fmt.Printf("%+v\n", item)
	})
}

func TestUtil(t *testing.T) {
	t.Log("Testing util...")

}

func TestNextId(t *testing.T) {

}

type request struct {
	BaseModel
	StartTime        time.Time   `json:"startTime" from:"startTime"`                         // 使用时间-开始时间
	EndTime          time.Time   `json:"endTime" from:"endTime"`                             // 使用时间-结束时间
	Type             int32       `json:"type,omitempty" from:"type"`                         // 优惠券类型 满减、指定商品等等
	Price            int32       `json:"price,omitempty" from:"price"`                       // 优惠金额分
	Full             int32       `json:"full,omitempty" from:"full"`                         // 满减条件(分)
	Name             string      `json:"name,omitempty" from:"name"`                         // 优惠券名称
	Remarks          string      `json:"remarks,omitempty" from:"remarks"`                   // 使用须知
	RuleDescription  string      `json:"ruleDescription,omitempty" from:"ruleDescription"`   // 规则说明
	ScopeDescription string      `json:"scopeDescription,omitempty" from:"scopeDescription"` // 范围说明
	TargetIds        interface{} `json:"targetIds,string,omitempty" from:"targetIds"`        // 商品范围
	TotalNum         int32       `json:"totalNum,omitempty" from:"totalNum"`                 // 发行总量
	UserLimit        int32       `json:"userLimit,omitempty" from:"userLimit"`               // 用户限领
	ReceiveStartTime time.Time   `json:"receiveStartTime" from:"receiveStartTime"`           // 领取时间-开始时间
	ReceiveEndTime   string      `json:"receiveEndTime" from:"receiveEndTime"`               // 领取时间-结束时间
	Disabled         bool        `json:"disabled,omitempty" from:"disabled"`                 // 禁用
	Sale             bool        `json:"sale,omitempty" from:"sale"`                         // 上架
	ValidityMinute   int32       `json:"validityMinute,omitempty" from:"validityMinute"`     // 领取后有效期分钟数
	SingleClaim      int32       `json:"singleClaim,omitempty" from:"singleClaim"`           // 单次领取次数
	Group            int32       `json:"group,omitempty" from:"group"`                       // 优惠券分组
}

func TestToSqlMap(t *testing.T) {
	r := request{
		Name:             "优惠券名称",
		Type:             1,
		Group:            1,
		ValidityMinute:   0,
		TotalNum:         1000,
		UserLimit:        1000,
		SingleClaim:      1,
		Price:            100,
		ReceiveStartTime: time.Now(),
		ReceiveEndTime:   "",
		StartTime:        time.Now(),
		EndTime:          time.Now(),
		Full:             100,
		RuleDescription:  "优惠券使用说明",
		ScopeDescription: "优惠券使用范围",
		Remarks:          "优惠券使用须知",
		Sale:             true,
		Disabled:         false,
		TargetIds:        NewInt8Array(1, 2, 3),
	}
	sqlMap := toSqlMap(context.Background(), r, map[string]interface{}{
		"Remarks": func(value interface{}) (interface{}, error) {
			return "使用须知x", errors.New("xxxx")
		},
		"RuleDescription": func(value interface{}) (interface{}, error) {
			return "使用须知x2", nil
		},
		"Disabled": func(value interface{}) (string, interface{}, error) {
			return "d2", ToMapIgnore, nil
		},
		"Full": ToMapIgnore,
		"id":   10,
		"id2": func() (interface{}, error) {
			return "ToMapIgnore", nil
		},
		"t2": 10,
		"TargetIds": func(value interface{}) (interface{}, error) {
			return ToMapIgnore, nil
		},
		"UserLimit": NewMapParse("x_UserLimit", func(value interface{}) (interface{}, error) {
			return "xxx", nil
		}),
	}, MapDateParse("ReceiveEndTime", time.DateTime, ToMapIgnore))
	for k, v := range sqlMap {
		t.Log(k, v)
	}
}

func TestAcv(t *testing.T) {
}

func TestFindIDs(t *testing.T) {

}

func TestFindIds(t *testing.T) {

	ids := FindIds[request]([]request{
		{
			TargetIds: NewInt8Array(1, 2, 3),
		},
		{
			TargetIds: NewInt8Array(2, 2, 4),
		},
	}, "TargetIds", true)
	t.Log(ids)
}

func TestSortByIds(t *testing.T) {
	ids := SortByIds([]request{
		{
			BaseModel: NewBaseModel(2),
		},
		{
			BaseModel: NewBaseModel(1),
		},
	}, []int64{2, 1, 2})
	for _, id := range ids {
		t.Log(id.ID)
	}
}

func TestFindLinksSet(t *testing.T) {
	FindLinksSet([]request{}, 2, func(t request) int64 {
		return t.ID
	}, func(t request) {

	})
}
