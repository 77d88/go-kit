package xdb

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/basic/xreflect"
	"github.com/77d88/go-kit/plugins/xlog"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	RegisterModels = make(map[string]GromModel)
	mu             sync.Mutex
	dbs            = make(map[string]*gorm.DB)
)

const dbStr string = "db"

// GetDB 获取数据库链接
func GetDB(name ...string) (*gorm.DB, error) {
	firstOrDefault := xarray.FirstOrDefault(name, dbStr)
	database, ok := dbs[firstOrDefault]
	if !ok {
		xlog.Errorf(nil, "数据库[%s]链接不存在", firstOrDefault)
		return nil, xerror.New("数据库链接不存在")
	}
	return database, nil
}

// AddModels 添加模型
func AddModels(dist ...GromModel) {
	mu.Lock()
	defer mu.Unlock()

	for _, v := range dist {
		if v == nil {
			continue
		}
		key := xreflect.Warp(v).InstPath()
		if _, ok := RegisterModels[key]; ok {
			continue
		}
		RegisterModels[key] = v
		xlog.Tracef(nil, "register model %s table %s", key, v.TableName())
	}
}

// AutoMigrateModel RegisterModels 自动迁移
func AutoMigrateModel(name ...string) error {
	db, err := GetDB(name...)
	if err != nil {
		return err
	}
	i := make([]interface{}, 0, len(RegisterModels))
	for _, v := range RegisterModels {
		i = append(i, v)
	}
	xlog.Tracef(nil, "auto migrate model %v", i)
	err = db.AutoMigrate(i...)
	if err != nil {
		xlog.Errorf(nil, "自动迁移失败: %+v", err)
	}
	return err
}

// KeyModel 模型主键
type KeyModel interface {
	GetID() int64
}

type Key struct {
	ID int64 `gorm:"comment:主键;primaryKey;autoIncrement:false" json:"id,string"` // 主键
}

type BaseModel struct {
	Key
	CreatedTime time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"createdTime"` // 创建时间
	UpdatedTime time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updatedTime"` // 更新时间
	DeletedTime gorm.DeletedAt `gorm:"comment:删除时间;index" json:"deletedTime"`          // 删除时间
}

func (b Key) GetID() int64 {
	return b.ID
}

func (b *Key) DefaultIgnoreUpdateFields() []string {
	return []string{}
}

func (b *Key) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = xid.NextId()
	}
	return
}

// GromModel 模型基础
type GromModel interface {
	KeyModel
	TableName() string
	DefaultIgnoreUpdateFields() []string
	InitData() []GromModel
}

func (b *BaseModel) DefaultIgnoreUpdateFields() []string {
	return []string{"created_time", "deleted_time"}
}

func (b *BaseModel) NewID() *BaseModel {
	b.ID = xid.NextId()
	return b
}

func (b *BaseModel) SetID(id int64) *BaseModel {
	b.ID = id
	return b
}
func (b *Key) InitData() []GromModel {
	return make([]GromModel, 0)
}

func NewId() Key {
	return Key{ID: NextId()}
}

func NewBaseModel(id ...int64) BaseModel {
	var i int64
	if len(id) == 0 || id[0] == 0 {
		i = NextId()
	} else {
		i = id[0]
	}
	return BaseModel{Key: Key{ID: i}}
}
