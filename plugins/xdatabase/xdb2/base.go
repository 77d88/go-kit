package xdb2

import (
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xlog"
)

var (
	mu        sync.Mutex
	dbs       = make(map[string]*DB)
	DefaultDB *DB
)

const DefaultDbLinkStr string = "db"

// GetDB 获取数据库链接
func GetDB(name ...string) (*DB, error) {
	if len(name) == 0 {
		return DefaultDB.clone(), nil
	}
	firstOrDefault := xarray.FirstOrDefault(name, DefaultDbLinkStr)
	database, ok := dbs[firstOrDefault]
	if !ok {
		xlog.Errorf(nil, "数据库[%s]链接不存在", firstOrDefault)
		return nil, xerror.New("数据库链接不存在")
	}
	return database.clone(), nil
}

type BaseModel struct {
	ID          int64     `gorm:"comment:主键;primaryKey;" json:"id,string"`                          // 主键
	CreatedTime time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdTime" db:"created_time"` // 创建时间
	UpdatedTime time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedTime" db:"updated_Time"` // 更新时间
	DeletedTime time.Time `gorm:"comment:删除时间;index" json:"deletedTime" db:"deleted_time"`          // 删除时间
}
