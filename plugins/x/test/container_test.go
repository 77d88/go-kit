package test

import (
	"testing"

	"github.com/77d88/go-kit/basic/xconfig/str_scanner"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"gorm.io/gorm"
)

func init() {
	str_scanner.Default(`{"server":{"port":9981,"debug":true},"db":{"dns":"host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai","logger":true},"redis":{"addr":"127.0.0.1:6666","pass":"test"}}`)
}

func TestBasic(t *testing.T) {
	x.Use("123")
	x.Use("1234", "config.str")
	get, err := x.Get[string]()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(get)

	get, err = x.Get[string]("config.str")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(get)
}

func TestConstructor(t *testing.T) {
	x.Use("34412")
	x.Use(func() (string, error) {
		return "123", nil
	}, "test")
	get, err := x.Get[string]("test")
	if err != nil {
		t.Error(err)
	}
	x.Use(xdb.NewEng)
	t.Log(get)

	for i := 0; i < 10; i++ {
		go func() {
			db, err2 := x.Get[*gorm.DB]()
			if err2 != nil {
				t.Error(err)
			}
			t.Logf("db: %v", db)
		}()
	}
	type a struct {
		DB   *gorm.DB
		Str  string `x:"name=test"`
		Str2 string
	}

	find, err := x.Find[a]()
	if err != nil {
		t.Error(err)
	}
	t.Logf("db2: %v", find)

	select {}
}
