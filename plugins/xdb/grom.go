package xdb

import (
	"context"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xreflect"
	"github.com/77d88/go-kit/plugins/xlog"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Init 创建pgSql链接
// host对应着远程postgresql数据库地址
// port为数据库端口,默认5432
// user为数据库用户,在安装的时候会创建postgres用户,也可使用自己创建的用户
// password为数据库用户对应的密码
// dbname为需要连接的具体实例库名

type DataSource struct {
	context.Context
	Config *Config
	DB     *gorm.DB
}

// New 创建链接 DataSource
func New(name ...string) *DataSource {
	db, err := GetDB(name...)
	if err != nil {
		return nil
	}
	return &DataSource{
		Context: context.Background(),
		DB:      db,
	}
}

func wrap(c context.Context, db *gorm.DB) *DataSource {
	return &DataSource{
		Context: c,
		DB:      db,
	}
}

func Ctx(c context.Context, names ...string) *DataSource {
	db := New(names...)
	return db.WithContext(c)
}
func (d *DataSource) WithContext(c context.Context) *DataSource {
	get := c.Value(CtxTransactionKey) // 有事务优先获取事务
	if get != nil {
		return get.(*DataSource)
	}
	d.Context = c
	d.DB = d.DB.WithContext(c)
	return d
}
func (d *DataSource) Session(session ...*gorm.Session) *DataSource {
	if session == nil || len(session) == 0 {
		session = append(session, &gorm.Session{})
	}
	if d == nil {
		return wrap(context.TODO(), d.DB.Session(session[0]))
	}
	return wrap(d.Context, d.DB.Session(session[0]))
}

func (d *DataSource) Table(name string, args ...interface{}) *DataSource {
	return wrap(d.Context, d.DB.Table(name, args...))
}

func (d *DataSource) unWrap() *gorm.DB {
	return d.DB
}

func (d *DataSource) Model(val interface{}) *DataSource {
	d.DB = d.DB.Model(val)
	return d
}

func (d *DataSource) Select(query interface{}, args ...interface{}) (tx *DataSource) {
	d.DB = d.DB.Select(query, args...)
	return d
}

func (d *DataSource) Where(query interface{}, args ...interface{}) *DataSource {
	d.DB = d.DB.Where(query, args...)
	return d
}

func (d *DataSource) WithId(ids ...int64) *DataSource {
	is := xarray.Filter(ids, func(i int, item int64) bool {
		return item > 0
	})
	if len(is) == 0 {
		return d
	}
	d.DB = d.DB.Where("id in ?", ids).Limit(len(ids))
	return d
}
func (d *DataSource) Order(value interface{}) *DataSource {
	d.DB = d.DB.Order(value)
	return d
}

func (d *DataSource) Having(query interface{}, args ...interface{}) *DataSource {
	d.DB = d.DB.Having(query, args)
	return d
}

func (d *DataSource) Group(name string) *DataSource {
	d.DB = d.DB.Group(name)
	return d
}

// XWhere 条件查询 条件判定为true才增加
func (d *DataSource) XWhere(condition bool, query interface{}, args ...interface{}) *DataSource {
	if !condition {
		return d
	}
	return d.Where(query, args...)
}

func (d *DataSource) Append(condition bool, f func(d *DataSource) *DataSource) *DataSource {
	if condition {
		d.DB = f(d).DB
		return d
	}
	return d
}

func (d *DataSource) Raw(sql string, values ...interface{}) *DataSource {
	d.DB = d.DB.Raw(sql, values...)
	return d
}

func (d *DataSource) Limit(limit int) *DataSource {
	d.DB = d.DB.Limit(limit)
	return d
}
func (d *DataSource) Offset(offset int) *DataSource {
	d.DB = d.DB.Offset(offset)
	return d
}
func (d *DataSource) IdDesc() *DataSource {
	d.DB = d.DB.Order("id desc")
	return d
}
func (d *DataSource) IdAsc() *DataSource {
	d.DB = d.DB.Order("id asc")
	return d
}

func (d *DataSource) PageSearch(page PageRequest) *DataSource {
	offset, limit := page.Limit()
	d.DB = d.DB.Offset(offset).Limit(limit)
	return d
}

func (d *DataSource) Page(page, size int) *DataSource {
	d.DB = d.DB.Offset((page - 1) * size).Limit(size)
	return d
}

func (d *DataSource) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *DataSource {
	d.DB = d.DB.Scopes(funcs...)
	return d
}
func (d *DataSource) Unscoped() *DataSource {
	d.DB = d.DB.Unscoped()
	return d
}
func (d *DataSource) Delete(value interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Delete(value, conds...))
}

func (d *DataSource) DeleteById(dest interface{}, id int64) *Result {
	if id <= 0 {
		return emptyResult
	}
	return warpResult(d.DB.Model(dest).Limit(1).Delete("id = ?", id))
}

func (d *DataSource) DeleteByIds(dest interface{}, ids ...int64) *Result {
	if len(ids) == 0 {
		return emptyResult
	}
	return warpResult(d.DB.Model(dest).Delete("id in (?)", ids))
}

// DeleteUnscoped 删除数据 物理删除
func (d *DataSource) DeleteUnscoped(value interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Unscoped().Delete(value, conds...))
}
func (d *DataSource) DeleteUnscopedById(dest interface{}, id int64) *Result {
	if id <= 0 {
		return emptyResult
	}
	return warpResult(d.DB.Unscoped().Model(dest).Limit(1).Delete("id = ?", id))
}
func (d *DataSource) DeleteUnscopedByIds(dest interface{}, ids ...int64) *Result {
	if len(ids) == 0 {
		return emptyResult
	}
	return warpResult(d.DB.Unscoped().Model(dest).Delete("id in (?)", ids))
}

func (d *DataSource) Find(dest interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Find(dest, conds...))
}
func (d *DataSource) FindByIds(dest interface{}, ids ...int64) *Result {
	if len(ids) == 0 {
		return emptyResult
	}
	return warpResult(d.DB.Find(dest, "id in (?)", ids))
}

// FindLinks 查询关联Id集合 fields 支持 int8Array 和 int64 其余不支持 默认去重查询 不保证排序
// source 可以为集合 也可以为单个对象
// to 为关联对象 实际数据库查询对象
// fields 为关联字段
// 示例：xdb.GetDB().FindLinks(&orders, &User{}, "UserID","SendId")
func (d *DataSource) FindLinks(source interface{}, to interface{}, fields ...string) *Result {
	// source 是否为集合
	ids := make([]int64, 0)
	for _, field := range fields {
		findIds := FindIds(xreflect.ToSlice(source), field, false)
		ids = append(ids, findIds...)
	}
	ids = xarray.Union(ids)

	if len(ids) == 0 {
		return emptyResult
	}
	return warpResult(d.DB.Find(to, "id in (?)", ids))
}

func (d *DataSource) FindById(dest interface{}, id int64) *Result {
	if id <= 0 {
		return emptyResult
	}
	return warpResult(d.DB.Find(dest, "id = ?", id))
}
func (d *DataSource) First(dest interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.First(dest, conds...))
}
func (d *DataSource) FirstById(dest interface{}, id int64) *Result {
	if id <= 0 {
		return &Result{
			Error: gorm.ErrRecordNotFound,
		}
	}
	return warpResult(d.DB.First(dest, "id = ?", id))
}
func (d *DataSource) Scan(dest interface{}) *Result {
	return warpResult(d.DB.Scan(dest))
}
func (d *DataSource) Count(count *int64) *Result {
	return warpResult(d.DB.Count(count))
}
func (d *DataSource) Create(value interface{}) *Result {
	return warpResult(d.DB.Create(value))
}
func (d *DataSource) Update(column string, value interface{}) *Result {
	return warpResult(d.DB.Update(column, value))
}
func (d *DataSource) Updates(values interface{}) *Result {
	return warpResult(d.DB.Updates(values))
}

// FindPage 分页查询
func (d *DataSource) FindPage(page PageRequest, result any, count *int64) *Result {
	if d.DB.Statement.Model == nil {
		d.Model(result)
	}
	offset, limit := page.Limit()
	if !page.IsNotCounted() && count != nil {
		// 首页统计总数
		if result := d.Count(count); result.Error != nil {
			return result
		}
		if *count <= int64(offset) {
			return emptyResult
		}
	}
	return d.IdDesc().Offset(offset).Limit(limit).Find(result)
}

func (d *DataSource) SaveMap(s GromModel, obj interface{}, mapping ...interface{}) *Result {
	m := toSqlMap(d, obj, mapping...)
	var id int64
	// 获取出ID 单独处理
	for k, v := range m {
		if strings.ToLower(k) == "id" {
			delete(m, k)
			if i, ok := v.(int64); ok {
				id = i
			}
		}
	}
	if id > 0 {
		result := d.Model(s).Where("id = ?", id).Updates(m)
		result.RowId = id
		return result
	} else {
		saveId := NextId()
		m["id"] = saveId
		m["created_time"] = time.Now()
		m["updated_time"] = time.Now()
		m["deleted_time"] = nil
		result := d.Model(s).Create(m)
		result.RowId = saveId
		return result
	}
}

func (d *DataSource) Exec(sql string, values ...interface{}) *Result {
	return warpResult(d.DB.Exec(sql, values...))
}

func (d *DataSource) Dispose() error {
	sqlDB, err := d.DB.DB() // 获取底层的 sql.DataSource 对象
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(d.Config.Dns, "password=******* ")
	xlog.Warnf(nil, "close db conn %s<%s>", d.Config.DbLinkName, maskedStr)
	return sqlDB.Close()
}
