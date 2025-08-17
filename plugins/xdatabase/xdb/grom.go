package xdb

import (
	"context"
	"regexp"

	"github.com/77d88/go-kit/plugins/xlog"

	"gorm.io/gorm"
)

func C(c context.Context, names ...string) *gorm.DB {
	return DB(names...).WithContext(c)
}

func G[T any](names ...string) gorm.Interface[T] {
	return gorm.G[T](DB(names...))
}

func DB(names ...string) *gorm.DB {
	db, err := GetDB(names...)
	if err != nil {
		xlog.Errorf(nil, "数据库[%v]链接不存在", names)
		return nil
	}
	return db
}

// New 创建pgSql链接
// host对应着远程postgresql数据库地址
// port为数据库端口,默认5432
// user为数据库用户,在安装的时候会创建postgres用户,也可使用自己创建的用户
// password为数据库用户对应的密码
// dbname为需要连接的具体实例库名

type db struct {
	DB     *gorm.DB
	Config *Config
}

//
//func (d *db) wrap(db *gorm.db) *db {
//	return &db{
//		db:     db,
//		Config: d.Config,
//	}
//}
//

//func (d *db) SQLDB() (*sql.db, error) {
//	return d.db.db()
//}
//func (d *db) WithCtx(c context.Context) *db {
//	get := GetCtxTran(c) // 有事务优先获取事务
//	if get != nil {
//		return get
//	}
//	return d.wrap(d.db.WithContext(c))
//}
//func (d *db) Session(session ...*gorm.Session) *db {
//	if session == nil || len(session) == 0 {
//		session = append(session, &gorm.Session{})
//	}
//	return d.wrap(d.db.Session(session[0]))
//}
//
//func (d *db) Table(name string, args ...interface{}) *db {
//	d.db = d.db.Table(name, args...)
//	return d
//}
//
//func (d *db) unWrap() *gorm.db {
//	return d.db
//}
//
//func (d *db) Model(val interface{}) *db {
//	d.db = d.db.Model(val)
//	return d
//}
//
//func (d *db) Select(query interface{}, args ...interface{}) (tx *db) {
//	d.db = d.db.Select(query, args...)
//	return d
//}
//
//func (d *db) Where(query interface{}, args ...interface{}) *db {
//	d.db = d.db.Where(query, args...)
//	return d
//}
//
//func (d *db) WithId(ids ...int64) *db {
//	is := xarray.Filter(ids, func(i int, item int64) bool {
//		return item > 0
//	})
//	if len(is) == 0 {
//		return d
//	}
//	if len(is) == 1 {
//		return d.Where("id = ?", is[0])
//	}
//	d.db = d.db.Where("id in ?", ids).Limit(len(ids))
//	return d
//}
//func (d *db) Order(value interface{}) *db {
//	d.db = d.db.Order(value)
//	return d
//}
//
//func (d *db) Having(query interface{}, args ...interface{}) *db {
//	d.db = d.db.Having(query, args)
//	return d
//}
//
//func (d *db) Group(name string) *db {
//	d.db = d.db.Group(name)
//	return d
//}
//
//// XWhere 条件查询 条件判定为true才增加
//func (d *db) XWhere(condition bool, query interface{}, args ...interface{}) *db {
//	if !condition {
//		return d
//	}
//	return d.Where(query, args...)
//}
//
//func (d *db) Append(condition bool, f func(d *db) *db) *db {
//	if condition {
//		d.db = f(d).db
//		return d
//	}
//	return d
//}
//
//func (d *db) Raw(sql string, values ...interface{}) *db {
//	d.db = d.db.Raw(sql, values...)
//	return d
//}
//
//func (d *db) Limit(limit int) *db {
//	d.db = d.db.Limit(limit)
//	return d
//}
//func (d *db) Offset(offset int) *db {
//	d.db = d.db.Offset(offset)
//	return d
//}
//func (d *db) IdDesc() *db {
//	d.db = d.db.Order("id desc")
//	return d
//}
//func (d *db) IdAsc() *db {
//	d.db = d.db.Order("id asc")
//	return d
//}
//
//func (d *db) PageSearch(page PageRequest) *db {
//	offset, limit := page.Limit()
//	return d.Offset(offset).Limit(limit)
//}
//
//func (d *db) Page(page, size int) *db {
//	return d.Offset((page - 1) * size).Limit(size)
//}
//
//func (d *db) Scopes(funcs ...func(*db) *db) *db {
//	gormFuncs := make([]func(*gorm.db) *gorm.db, len(funcs))
//	for i, fn := range funcs {
//		f := fn
//		gormFuncs[i] = func(gdb *gorm.db) *gorm.db {
//			return f(d.wrap(gdb)).unWrap()
//		}
//	}
//	d.db = d.db.Scopes(gormFuncs...)
//	return d
//}
//func (d *db) Unscoped() *db {
//	d.db = d.db.Unscoped()
//	return d
//}
//func (d *db) Delete(value interface{}, conds ...interface{}) *Result {
//	return warpResult(d.db.Delete(value, conds...))
//}
//
//func (d *db) DeleteById(dest interface{}, ids ...int64) *Result {
//	if len(ids) == 0 {
//		return emptyResult
//	}
//	if len(ids) == 1 {
//		return warpResult(d.db.Model(dest).Limit(1).Delete("id = ?", ids[0]))
//	}
//	return warpResult(d.db.Model(dest).Limit(len(ids)).Delete("id in (?)", ids))
//}
//
//// DeleteUnscoped 删除数据 物理删除
//func (d *db) DeleteUnscoped(value interface{}, conds ...interface{}) *Result {
//	return warpResult(d.db.Unscoped().Delete(value, conds...))
//}
//func (d *db) DeleteUnscopedByIds(dest interface{}, ids ...int64) *Result {
//	if len(ids) == 0 {
//		return emptyResult
//	}
//	if len(ids) == 1 {
//		return warpResult(d.db.Unscoped().Model(dest).Delete("id = ?", ids[0]))
//	}
//	return warpResult(d.db.Unscoped().Model(dest).Delete("id in (?)", ids).Limit(len(ids)))
//}
//
//func (d *db) Find(dest interface{}, conds ...interface{}) *Result {
//	return warpResult(d.db.Find(dest, conds...))
//}
//func (d *db) FindByIds(dest interface{}, ids ...int64) *Result {
//	if len(ids) == 0 {
//		return emptyResult
//	}
//	return warpResult(d.db.Find(dest, "id in (?)", ids))
//}
//
//// FindLinks 查询关联Id集合 fields 支持 int8Array 和 int64 其余不支持 默认去重查询 不保证排序
//// source 可以为集合 也可以为单个对象
//// to 为关联对象 实际数据库查询对象
//// fields 为关联字段
//// 示例：xdb.GetDB().FindLinks(&orders, &User{}, "UserID","SendId")
//func (d *db) FindLinks(source interface{}, to interface{}, fields ...string) *Result {
//	// source 是否为集合
//	ids := make([]int64, 0)
//	for _, field := range fields {
//		findIds := FindIds(xreflect.ToSlice(source), field, false)
//		ids = append(ids, findIds...)
//	}
//	ids = xarray.Union(ids)
//
//	if len(ids) == 0 {
//		return emptyResult
//	}
//	return warpResult(d.db.Find(to, "id in (?)", ids))
//}
//func (d *db) Take(dest interface{}) *Result {
//	return warpResult(d.db.Take(dest))
//}
//func (d *db) First(dest interface{}, conds ...interface{}) *Result {
//	return warpResult(d.db.Take(dest, conds...))
//}
//func (d *db) Scan(dest interface{}) *Result {
//	return warpResult(d.db.Scan(dest))
//}
//func (d *db) Count(count *int64) *Result {
//	return warpResult(d.db.Count(count))
//}
//func (d *db) Create(value interface{}) *Result {
//	return warpResult(d.db.Create(value))
//}
//func (d *db) Update(column string, value interface{}) *Result {
//	return warpResult(d.db.Update(column, value))
//}
//func (d *db) Updates(values interface{}) *Result {
//	return warpResult(d.db.Updates(values))
//}
//
//// FindPage 分页查询
//func (d *db) FindPage(page PageRequest, result any, count *int64) *Result {
//	if d.db.Statement.Model == nil {
//		d.Model(result)
//	}
//	offset, limit := page.Limit()
//	if !page.IsNotCounted() && count != nil {
//		// 首页统计总数
//		if result := d.Count(count); result.Error != nil {
//			return result
//		}
//		if *count <= int64(offset) {
//			return emptyResult
//		}
//	}
//	return d.IdDesc().Offset(offset).Limit(limit).Find(result)
//}
//
//func (d *db) SaveMap(s GromModel, obj interface{}, mapping ...interface{}) *Result {
//	m := toSqlMap(obj, mapping...)
//	mdb := d
//	if s != nil {
//		mdb = d.Model(s)
//	}
//	var id int64
//	// 获取出ID 单独处理
//	for k, v := range m {
//		if strings.ToLower(k) == "id" {
//			delete(m, k)
//			if i, ok := v.(int64); ok {
//				id = i
//			}
//		}
//	}
//	if id > 0 {
//		result := mdb.Where("id = ?", id).Updates(m)
//		result.RowId = id
//		return result
//	} else {
//		saveId := NextId()
//		m["id"] = saveId
//		m["created_time"] = time.Now()
//		m["updated_time"] = time.Now()
//		m["deleted_time"] = nil
//		result := mdb.Create(m)
//		result.RowId = saveId
//		return result
//	}
//}
//
//func (d *db) Exec(sql string, values ...interface{}) *Result {
//	return warpResult(d.db.Exec(sql, values...))
//}

func (d *db) Dispose() error {
	sqlDB, err := d.DB.DB() // 获取底层的 sql.db 对象
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(d.Config.Dns, "password=******* ")
	xlog.Warnf(nil, "close db conn %s<%s>", d.Config.DbLinkName, maskedStr)
	return sqlDB.Close()
}
