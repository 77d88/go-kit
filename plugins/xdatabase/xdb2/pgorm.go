package xdb2

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 是 pgxorm 的主对象，模仿 GORM 的 *gorm.DB
type DB struct {
	pool *pgxpool.Pool
}

// Logger 定义日志接口，可对接 zap、logrus 等
type Logger interface {
	Print(v ...interface{})
}

func C(c context.Context, name ...string) *Inst {
	db, err := GetDB(name...)
	if err != nil {
		return nil
	}
	return &Inst{
		pool: db.pool,
		cond: sq.Select().PlaceholderFormat(sq.Dollar),
		ctx:  c,
	}
}

// NewDB 创建一个新的 DB 实例
//func NewDB(pool *pgxpool.Pool, opts ...Option) *DB {
//	db := &DB{
//		pool: pool,
//		cond: sqlbuilder.NewSelectBuilder(),
//	}
//	for _, opt := range opts {
//		opt(db)
//	}
//	return db
//}

// Option 用于配置 DB
// type Option func(*DB)
//
// // WithLogger 设置日志器
//
//	func WithLogger(logger Logger) Option {
//		return func(db *DB) {
//			db.logger = logger
//		}
//	}
//
// // WithContext 设置上下文
//
//	func WithContext(ctx context.Context) Option {
//		return func(db *DB) {
//			db.ctx = ctx
//		}
//	}
//
// // Model 设置操作的模型（结构体）
//
//	func (db *DB) Model(model interface{}) *DB {
//		newDB := db.clone()
//		newDB.model = parseModel(model)
//		newDB.cond.From(newDB.model.TableName)
//		return newDB
//	}
//
// // Where 添加 WHERE 条件
//
//	func (db *DB) Where(query interface{}, args ...interface{}) *DB {
//		newDB := db.clone()
//		switch q := query.(type) {
//		case string:
//			newDB.cond.Where(q, args...)
//		case map[string]interface{}:
//			conds := []string{}
//			var a []interface{}
//			for k, v := range q {
//				conds = append(conds, fmt.Sprintf("%s = ?", k))
//				a = append(a, v)
//			}
//			newDB.cond.Where(strings.Join(conds, " AND "), a...)
//		}
//		return newDB
//	}
//
// // Limit 设置 LIMIT
//
//	func (db *DB) Limit(limit int) *DB {
//		newDB := db.clone()
//		newDB.cond.Limit(limit)
//		return newDB
//	}
//
// // Offset 设置 OFFSET
//
//	func (db *DB) Offset(offset int) *DB {
//		newDB := db.clone()
//		newDB.cond.Offset(offset)
//		return newDB
//	}
//
// // Order 设置排序
//
//	func (db *DB) Order(order string) *DB {
//		newDB := db.clone()
//		newDB.cond.OrderBy(order)
//		return newDB
//	}
//
// // Find 查询多条记录
//
//	func (db *DB) Find(dest interface{}) error {
//		if db.model == nil {
//			return errors.New("model is not set, use Model()")
//		}
//
//		// 构建 SQL
//		sql, args := db.cond.Select(db.model.SelectFields()...).MustBuild()
//
//		rows, err := db.query(sql, args...)
//		if err != nil {
//			return err
//		}
//		defer rows.Close()
//
//		sliceValue := reflect.ValueOf(dest)
//		if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
//			return errors.New("dest must be a pointer to slice")
//		}
//
//		elemType := sliceValue.Elem().Type().Elem()
//		for rows.Next() {
//			elem := reflect.New(elemType).Interface()
//			if err := db.scanIntoStruct(elem, rows); err != nil {
//				return err
//			}
//			sliceValue.Elem().Set(reflect.Append(sliceValue.Elem(), reflect.ValueOf(elem).Elem()))
//		}
//
//		return rows.Err()
//	}
//
// // First 查询第一条
//
//	func (db *DB) First(dest interface{}) error {
//		return db.Limit(1).Find(dest)
//	}
//
// // Create 插入一条记录
//
//	func (db *DB) Create(value interface{}) error {
//		if db.model == nil {
//			db.model = parseModel(value)
//		}
//
//		fields, placeholders, values := db.model.ExtractValues(value)
//		builder := sqlbuilder.InsertInto(db.model.TableName).Columns(fields...).Values(placeholders...)
//
//		sql, args := builder.MustBuild()
//
//		_, err := db.exec(sql, args...)
//		return err
//	}
//
// // Save 更新或插入（根据主键判断）
//
//	func (db *DB) Save(value interface{}) error {
//		if db.model == nil {
//			db.model = parseModel(value)
//		}
//
//		idValue := db.model.GetIDValue(value)
//		if idValue == nil {
//			return db.Create(value)
//		}
//
//		// 尝试更新
//		fields, values := db.model.ExtractUpdateFields(value)
//		builder := sqlbuilder.Update(db.model.TableName)
//		for i, f := range fields {
//			builder.Set(f, values[i])
//		}
//		builder.Where(fmt.Sprintf("%s = ?", db.model.PrimaryKey), idValue)
//
//		sql, args := builder.MustBuild()
//		tag, err := db.exec(sql, args...)
//		if err != nil {
//			return err
//		}
//
//		// 如果未更新，则插入
//		if tag.RowsAffected() == 0 {
//			return db.Create(value)
//		}
//		return nil
//	}
//
// // Delete 删除
//
//	func (db *DB) Delete(value interface{}, conds ...interface{}) error {
//		if db.model == nil {
//			db.model = parseModel(value)
//		}
//
//		builder := sqlbuilder.DeleteFrom(db.model.TableName)
//
//		if len(conds) > 0 {
//			if whereStr, ok := conds[0].(string); ok {
//				builder.Where(whereStr, conds[1:]...)
//			}
//		} else {
//			idValue := db.model.GetIDValue(value)
//			if idValue != nil {
//				builder.Where(fmt.Sprintf("%s = ?", db.model.PrimaryKey), idValue)
//			}
//		}
//
//		sql, args := builder.MustBuild()
//		_, err := db.exec(sql, args...)
//		return err
//	}
//
// // Exec 原生执行
//
//	func (db *DB) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
//		if db.tx != nil {
//			return db.tx.Exec(db.ctx, sql, args...)
//		}
//		return db.pool.Exec(db.ctx, sql, args...)
//	}
//
// // query 执行查询
//
//	func (db *DB) query(sql string, args ...interface{}) (pgx.Rows, error) {
//		if db.tx != nil {
//			return db.tx.Query(db.ctx, sql, args...)
//		}
//		return db.pool.Query(db.ctx, sql, args...)
//	}
//
// // exec 执行语句
//
//	func (db *DB) exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
//		return db.Exec(sql, args...)
//	}
//
// // scanIntoStruct 将 pgx.Rows 扫描到结构体
//
//	func (db *DB) scanIntoStruct(dest interface{}, rows pgx.Rows) error {
//		columns := rows.FieldDescriptions()
//		values := make([]interface{}, len(columns))
//
//		for i := range values {
//			values[i] = new(interface{})
//		}
//
//		err := rows.Scan(values...)
//		if err != nil {
//			return err
//		}
//
//		return db.model.MapScan(dest, values)
//	}
//
// // clone 复制当前 DB 实例（用于链式调用）
func (db *DB) clone() *DB {
	return &DB{
		pool: db.pool,
	}
}

//
//// Error 返回当前错误
//func (db *DB) Error() error {
//	return db.error
//}
//
//// Rows 返回原始 rows（可选）
//func (db *DB) Rows() (pgx.Rows, error) {
//	sql, args := db.cond.MustBuild()
//	return db.query(sql, args...)
//}
