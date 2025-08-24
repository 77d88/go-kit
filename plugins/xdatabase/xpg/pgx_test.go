package xpg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/77d88/go-kit/basic/xtime"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	BaseModel
	UpdateUser          int64    `gorm:"comment:更新人"`
	Password            string   `gorm:"comment:后台登录密码" json:"password"`   // 后台登录密码
	Disabled            bool     `gorm:"comment:是否禁用" json:"disabled"`       // 是否禁用
	Username            string   `gorm:"comment:后台登录名称" json:"username"`   // 后台登录名称
	Nickname            string   `gorm:"comment:后台显示名称" json:"nickname"`   // 后台显示名称
	Avatar              []int64  `gorm:"comment:头像" json:"avatar"`             // 头像
	Roles               []int64  `gorm:"comment:系统角色" json:"roles"`          // 系统角色
	Permission          []int64  `gorm:"comment:系统独立权限" json:"permission"` // 系统独立权限
	Email               string   `gorm:"comment:邮箱" json:"email"`
	IsReLogin           bool     `gorm:"comment:是否需要重新登录" json:"isReLogin"`
	ReLoginDesc         string   `gorm:"comment:重新登录描述" json:"reLoginDesc"`
	PermissionCodes     []string `gorm:"comment:权限码" json:"permissionCodes"`     // 冗余 集合Permission里面的所有
	RolePermissionCodes []string `gorm:"comment:角色码" json:"RolePermissionCodes"` // 冗余 集合Roles里面的所有Permission Code
	_codes              []string `db:"-"`                                           // 本地计算的code
	_isCalcCodes        bool     `db:"-"`                                           // 是否计算
}

// TableName Res's table name
func (*User) TableName() string {
	return "s_sys_user"
}
func init() {
	pool, err := pgxpool.New(context.Background(), "host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable")
	if err != nil {
		panic(err)
	}
	DefaultDB = &DB{pool: pool}
}

func TestName(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	// 预热
	var result int
	err = pool.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		panic(err)
	}

	// 多次测试
	for i := 0; i < 100; i++ {
		start := time.Now()
		err := pool.QueryRow(context.Background(), "SELECT count(*) FROM \"s_user\" WHERE \"s_user\".\"deleted_time\" IS NULL").Scan(&result)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		execTime := time.Since(start)
		fmt.Printf("Attempt %d: %v ms result %d\n", i+1, execTime.Milliseconds(), result)
	}
}

func TestSqlbuilder(t *testing.T) {
	sb := sq.Select("name", "age").From("user")
	sb = sb.Where("b && ?", "123").Columns("name2", "age")
	fmt.Println(sb.ToSql())
}

func TestPGorm(t *testing.T) {

	type User struct {
		BaseModel
		Username string `db:"username"`
	}
	for i := 0; i < 10; i++ {
		var result []User
		c := C(context.Background())
		inv := xtime.NewTimeInterval()
		r := c.Table("s_user").Debug().Find(&result)
		t.Logf("result[%d:ms]: %+v:%+v ", inv.IntervalMs(), r, result)
	}
}
func TestConr(t *testing.T) {
	c := C(context.Background()).Model(&User{})
	inst := c.Copy()
	c = c.Where("id=?", 1)
	sql, i, err := c.ToSql()
	t.Logf("sql:%s, args:%v, err:%v", sql, i, err)
	c = inst.Where("id=?", 2)
	sql, i, err = c.ToSql()
	t.Logf("sql:%s, args:%v, err:%v", sql, i, err)

}

func TestModel(t *testing.T) {
	c := C(context.Background())
	var model []User
	r := c.Debug().Find(&model)
	t.Logf("result:%+v:%+v", r, model)
}

func TestUpdate(t *testing.T) {
	for i := 0; i < 10; i++ {
		c := C(context.Background()).Debug()
		result := c.Model(&User{}).Where("id=?", 1).Updates(map[string]interface{}{
			"disabled":    true,
			"update_user": -1,
		})
		t.Logf("result:%+v", result)
	}
}

func TestCreate(t *testing.T) {
	c := C(context.Background()).Debug()
	result := c.Create(&User{
		Disabled:   true,
		UpdateUser: -1,
	})
	t.Logf("result:%+v", result)
}

func TestSave(t *testing.T) {
	c := C(context.Background()).Debug()
	result := c.Save(&User{
		Disabled:   true,
		UpdateUser: -1,
		BaseModel:  BaseModel{ID: 910309449912389},
	}, func(m map[string]interface{}) {
		m["update_user"] = -2
	})
	t.Logf("result:%+v", result)
}

func TestFirst(t *testing.T) {
	c := C(context.Background()).Debug()
	result := &User{}
	r := c.Where("id=?", 999).First(result)
	t.Logf("result:%+v ==>%+v  %v", r, result, r.IsNotFound())
}

func TestCount(t *testing.T) {
	var count int64
	c := C(context.Background()).Debug()
	result := c.Model(&User{}).Count(&count)

	t.Logf("result:%+v ==>%d", result, count)
}

func TestCount2(t *testing.T) {
	var count string
	c := C(context.Background()).Debug()
	result := c.Raw("select count(1) from s_sys_user").Scan(&count)

	t.Logf("result:%+v ==>%s", result, count)
}

func TestFindPage(t *testing.T) {
	c := C(context.Background()).Debug()
	var model []User
	result := c.Model(&User{}).Where("id>?", 0).FindPage(&model, PageSearch{Page: 3, Size: 1}, true)
	t.Logf("result:%+v ==>%+v", result, model)
}
