package xdb2

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/77d88/go-kit/basic/xtime"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	pool, err := pgxpool.New(context.Background(), "host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	DefaultDB = &DB{pool: pool}
	type User struct {
		Id int64 `db:"id"`
	}
	for i := 0; i < 10; i++ {
		var result []User
		c := C(context.Background())

		inv := xtime.NewTimeInterval()
		err = c.Table("s_user").Select("*").Find(&result)
		t.Logf("[%d]result: %+v err: %v  %v", inv.IntervalMs(), result, err, errors.Is(err, pgx.ErrNoRows))
	}
}
