package xdb

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/xlog"
	"reflect"
	"runtime"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	RegisterModels = make(map[string]map[string]GromModel)
	mu             sync.Mutex
	dbs            = make(map[string]*gorm.DB)
	DefaultDB      *gorm.DB
)

const DefaultDbLinkStr string = "db"

// GetDB 获取数据库链接
func GetDB(name ...string) (*gorm.DB, error) {
	if len(name) == 0 {
		return DefaultDB, nil
	}
	firstOrDefault := xarray.FirstOrDefault(name, DefaultDbLinkStr)
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
		var dbname = DefaultDbLinkStr

		if v, ok := v.(DBNamer); ok {
			dbname = v.DbName()
		}
		re := RegisterModels[dbname]
		if re == nil {
			re = make(map[string]GromModel)
		}

		// 获取v的文件完整路径
		modelType := reflect.TypeOf(v)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		//var fileName string
		//if modelType != nil {
		//	// 尝试获取TableName方法的定义文件
		//	if tableNameMethod, found := modelType.MethodByName("TableName"); found {
		//		if tableNameMethod.Func.IsValid() {
		//			file, _ := runtime.FuncForPC(tableNameMethod.Func.Pointer()).FileLine(0)
		//			fileName = file
		//		} else {
		//			fileName = modelType.PkgPath()
		//		}
		//	} else {
		//		fileName = modelType.PkgPath()
		//	}
		//} else {
		//	// 如果无法获取模型信息，使用调用方文件
		//	_, file, _, _ := runtime.Caller(1)
		//	fileName = file
		//}
		_, file, line, _ := runtime.Caller(1)
		if re[v.TableName()] == nil {
			re[v.TableName()] = v
			xlog.Debugf(nil, "register %s => table %s(%s:%d)", dbname, v.TableName(), file, line)
		} else {
			xlog.Warnf(nil, "table %s(%s) => [%s:%d] already exists", dbname, v.TableName(), file, line)
		}
		RegisterModels[dbname] = re
	}
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

func (b *Key) GetID() int64 {
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

type DBNamer interface {
	DbName() string
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
