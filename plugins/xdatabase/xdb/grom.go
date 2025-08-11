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

type DB struct {
	DB *gorm.DB
	Config *Config
}

// New 创建链接 DB
func New(name ...string) *DB {
	db, err := GetDB(name...)
	if err != nil {
		return nil
	}
	return &DB{
		DB: db,
	}
}

func wrap(db *gorm.DB) *DB {
	return &DB{
		DB: db,
	}
}

func Ctx(c context.Context, names ...string) *DB {
	db := New(names...)
	return db.WithCtx(c)
}
func (d *DB) WithCtx(c context.Context) *DB {
	get := GetCtxTran(c) // 有事务优先获取事务
	if get != nil {
		return get
	}
	return wrap(d.DB.WithContext(c))
}
func (d *DB) Session(session ...*gorm.Session) *DB {
	if session == nil || len(session) == 0 {
		session = append(session, &gorm.Session{})
	}
	if d == nil {
		return wrap(d.DB.Session(session[0]))
	}
	return wrap(d.DB.Session(session[0]))
}

func (d *DB) Table(name string, args ...interface{}) *DB {
	d.DB = d.DB.Table(name, args...)
	return d
}

func (d *DB) unWrap() *gorm.DB {
	return d.DB
}

func (d *DB) Model(val interface{}) *DB {
	d.DB = d.DB.Model(val)
	return d
}

func (d *DB) Select(query interface{}, args ...interface{}) (tx *DB) {
	d.DB = d.DB.Select(query, args...)
	return d
}

func (d *DB) Where(query interface{}, args ...interface{}) *DB {
	d.DB = d.DB.Where(query, args...)
	return d
}

func (d *DB) WithId(ids ...int64) *DB {
	is := xarray.Filter(ids, func(i int, item int64) bool {
		return item > 0
	})
	if len(is) == 0 {
		return d
	}
	if len(is) == 1 {
		return d.Where("id = ?", is[0])
	}
	d.DB = d.DB.Where("id in ?", ids).Limit(len(ids))
	return d
}
func (d *DB) Order(value interface{}) *DB {
	d.DB = d.DB.Order(value)
	return d
}

func (d *DB) Having(query interface{}, args ...interface{}) *DB {
	d.DB = d.DB.Having(query, args)
	return d
}

func (d *DB) Group(name string) *DB {
	d.DB = d.DB.Group(name)
	return d
}

// XWhere 条件查询 条件判定为true才增加
func (d *DB) XWhere(condition bool, query interface{}, args ...interface{}) *DB {
	if !condition {
		return d
	}
	return d.Where(query, args...)
}

func (d *DB) Append(condition bool, f func(d *DB) *DB) *DB {
	if condition {
		d.DB = f(d).DB
		return d
	}
	return d
}

func (d *DB) Raw(sql string, values ...interface{}) *DB {
	d.DB = d.DB.Raw(sql, values...)
	return d
}

func (d *DB) Limit(limit int) *DB {
	d.DB = d.DB.Limit(limit)
	return d
}
func (d *DB) Offset(offset int) *DB {
	d.DB = d.DB.Offset(offset)
	return d
}
func (d *DB) IdDesc() *DB {
	d.DB = d.DB.Order("id desc")
	return d
}
func (d *DB) IdAsc() *DB {
	d.DB = d.DB.Order("id asc")
	return d
}

func (d *DB) PageSearch(page PageRequest) *DB {
	offset, limit := page.Limit()
	return d.Offset(offset).Limit(limit)
}

func (d *DB) Page(page, size int) *DB {
	return d.Offset((page - 1) * size).Limit(size)
}

func (d *DB) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *DB {
	d.DB = d.DB.Scopes(funcs...)
	return d
}
func (d *DB) Unscoped() *DB {
	d.DB = d.DB.Unscoped()
	return d
}
func (d *DB) Delete(value interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Delete(value, conds...))
}

func (d *DB) DeleteById(dest interface{}, ids ...int64) *Result {
	if len(ids) == 0 {
		return emptyResult
	}
	if len(ids) == 1 {
		return warpResult(d.DB.Model(dest).Limit(1).Delete("id = ?", ids[0]))
	}
	return warpResult(d.DB.Model(dest).Limit(len(ids)).Delete("id in (?)", ids))
}

// DeleteUnscoped 删除数据 物理删除
func (d *DB) DeleteUnscoped(value interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Unscoped().Delete(value, conds...))
}
func (d *DB) DeleteUnscopedByIds(dest interface{}, ids ...int64) *Result {
	if len(ids) == 0 {
		return emptyResult
	}
	if len(ids) == 1 {
		return warpResult(d.DB.Unscoped().Model(dest).Delete("id = ?", ids[0]))
	}
	return warpResult(d.DB.Unscoped().Model(dest).Delete("id in (?)", ids).Limit(len(ids)))
}

func (d *DB) Find(dest interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Find(dest, conds...))
}
func (d *DB) FindByIds(dest interface{}, ids ...int64) *Result {
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
func (d *DB) FindLinks(source interface{}, to interface{}, fields ...string) *Result {
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
func (d *DB) Take(dest interface{}) *Result  {
	return warpResult(d.DB.Take(dest))
}
func (d *DB) First(dest interface{}, conds ...interface{}) *Result {
	return warpResult(d.DB.Take(dest, conds...))
}
func (d *DB) Scan(dest interface{}) *Result {
	return warpResult(d.DB.Scan(dest))
}
func (d *DB) Count(count *int64) *Result {
	return warpResult(d.DB.Count(count))
}
func (d *DB) Create(value interface{}) *Result {
	return warpResult(d.DB.Create(value))
}
func (d *DB) Update(column string, value interface{}) *Result {
	return warpResult(d.DB.Update(column, value))
}
func (d *DB) Updates(values interface{}) *Result {
	return warpResult(d.DB.Updates(values))
}

// FindPage 分页查询
func (d *DB) FindPage(page PageRequest, result any, count *int64) *Result {
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

func (d *DB) SaveMap(s GromModel, obj interface{}, mapping ...interface{}) *Result {
	m := toSqlMap(obj, mapping...)
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

func (d *DB) Exec(sql string, values ...interface{}) *Result {
	return warpResult(d.DB.Exec(sql, values...))
}

func (d *DB) Dispose() error {
	sqlDB, err := d.DB.DB() // 获取底层的 sql.DB 对象
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(d.Config.Dns, "password=******* ")
	xlog.Warnf(nil, "close db conn %s<%s>", d.Config.DbLinkName, maskedStr)
	return sqlDB.Close()
}
