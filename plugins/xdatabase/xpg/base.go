package xpg

import (
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
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

type Model interface {
	TableName() string
}

type BaseModel struct {
	ID          int64     `json:"id,string" db:"id"`             // 主键
	CreatedTime time.Time `json:"createdTime" db:"created_time"` // 创建时间
	UpdatedTime time.Time `json:"updatedTime" db:"updated_time"` // 更新时间
	DeletedTime time.Time `json:"deletedTime" db:"deleted_time"` // 删除时间
}

func NewBaseModel(id ...int64) BaseModel {
	var nextID int64
	if len(id) == 0 || id[0] == 0 {
		nextID = xid.NextId()
	} else {
		nextID = id[0]
	}
	return BaseModel{
		ID:          nextID,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}
