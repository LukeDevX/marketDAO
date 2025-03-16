// 封装了DB操作CURD基本方法，用于被其它数据模型继承
package model

import (
	"RESTful-API/internal/constants" // 导入内部常量包
	"RESTful-API/internal/errno"     // 导入内部错误包
	"reflect"                        // 导入反射包，用于类型检查
	"strings"                        // 导入字符串处理包

	"fmt" // 导入格式化包

	"gorm.io/gorm" // 导入GORM包，用于数据库操作
)

var OrmMap map[string]*gorm.DB // 全局变量，存储不同名称的GORM DB实例

type BaseModel struct {
	db        *gorm.DB // GORM DB实例
	tableName string   // 数据库表名
	T         string   // 类型字段，未使用
}

// NewOrm 初始化并返回一个GORM DB实例
func NewOrm(alias ...string) *gorm.DB {
	ormName := "default" // 默认ORM名称
	if len(alias) > 0 {
		ormName = alias[0] // 如果传入别名，则使用别名
	}

	if OrmMap == nil {
		panic("[NewOrm] orm not init") // 如果OrmMap未初始化，抛出异常
	}

	db, ok := OrmMap[ormName] // 从OrmMap中获取对应的DB实例
	if !ok {
		panic("[NewOrm] orm not exists") // 如果不存在，抛出异常
	}

	return db // 返回DB实例
}

// Create 基础模型的插入方法
func (m *BaseModel) Create(dao interface{}) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(dao) // 获取dao的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果dao不是指针或为空，返回错误
	}
	db = db.Create(dao) // 执行插入操作
	if db.Error != nil {
		return db.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// Delete 基础模型的删除方法
func (m *BaseModel) Delete(filters map[string]interface{}, dao interface{}) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(dao) // 获取dao的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果dao不是指针或为空，返回错误
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Delete(dao) // 执行删除操作
	if db.Error != nil {
		return db.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// QueryOne 基础模型的查询单行方法
func (m *BaseModel) QueryOne(filters map[string]interface{}, dao interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(dao) // 获取dao的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果dao不是指针或为空，返回错误
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	err := db.Limit(1).Find(dao).Error // 执行查询操作，限制返回1条记录
	if OrmErr(err) != nil {
		return err // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// Update 基础模型的更新方法
func (m *BaseModel) Update(data map[string]interface{}, filters map[string]interface{}, limitNum ...int) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil || db == nil ||
		len(filters) <= 0 ||
		len(data) <= 0 {
		return 0, errno.InvalidParamError // 如果BaseModel、DB实例、过滤条件或数据为空，返回参数错误
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}

	if len(limitNum) > constants.DefaultZero {
		db = db.Limit(limitNum[0]).Updates(data) // 如果有数量限制，则添加限制并执行更新
	} else {
		db = db.Updates(data) // 否则直接执行更新
	}

	if db.Error != nil {
		return 0, db.Error // 如果出错，返回错误
	}
	return db.RowsAffected, nil // 返回受影响的行数和nil
}

// ListAndTotal 基础模型的分页查询方法，返回列表和总数
func (m *BaseModel) ListAndTotal(filters map[string]interface{}, page, pageSize int, list interface{}, orderBy ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(list) // 获取list的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errno.NonPointerOrEmptyError // 如果list不是指针或为空，返回错误
	}

	var total int64 // 用于存储总数

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Count(&total) // 计算总数
	if OrmErr(db.Error) != nil {
		return 0, db.Error // 如果出错，返回错误
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(list) // 执行分页查询
	if OrmErr(db.Error) != nil {
		return 0, db.Error // 如果出错，返回错误
	}
	return total, nil // 返回总数和nil
}

// TotalByJoin 通过JOIN查询总数
func (m *BaseModel) TotalByJoin(filters map[string]interface{}, joinSQL string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName).Joins(joinSQL) // 否则初始化一个新的DB实例并指定表名和JOIN条件
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}
	var total int64 // 用于存储总数
	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Count(&total) // 计算总数
	if OrmErr(db.Error) != nil {
		return 0, db.Error // 如果出错，返回错误
	}

	return total, nil // 返回总数和nil
}

// ListAndTotalByJoin 通过JOIN查询列表和总数
func (m *BaseModel) ListAndTotalByJoin(filters map[string]interface{}, fields []string, joinSQL string, page, pageSize int, list interface{}, orderBy ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(list) // 获取list的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errno.NonPointerOrEmptyError // 如果list不是指针或为空，返回错误
	}

	var total int64 // 用于存储总数

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Count(&total) // 计算总数
	if OrmErr(db.Error) != nil {
		return 0, db.Error // 如果出错，返回错误
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Select(fields).Joins(joinSQL).Find(list) // 执行分页查询并指定字段和JOIN条件
	if OrmErr(db.Error) != nil {
		return 0, db.Error // 如果出错，返回错误
	}
	return total, nil // 返回总数和nil
}

// List 基础模型的分页查询方法，返回列表
func (m *BaseModel) List(filters map[string]interface{}, page, pageSize int, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(list) // 获取list的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果list不是指针或为空，返回错误
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(list) // 执行分页查询
	if OrmErr(db.Error) != nil {
		return db.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// ListWithJoin 通过JOIN查询列表
func (m *BaseModel) ListWithJoin(filters map[string]interface{}, fields []string, joinSQL string, page, pageSize int, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(list) // 获取list的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果list不是指针或为空，返回错误
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Select(fields).Joins(joinSQL).Find(list) // 执行分页查询并指定字段和JOIN条件
	if OrmErr(db.Error) != nil {
		return db.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// ListNoPage 基础模型的查询方法，不进行分页
func (m *BaseModel) ListNoPage(filters map[string]interface{}, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(list) // 获取list的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果list不是指针或为空，返回错误
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value) // 添加过滤条件
		}
	}

	db = db.Find(list) // 执行查询操作
	if OrmErr(db.Error) != nil {
		return db.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// Sum 基础模型的求和方法
func (m *BaseModel) Sum(field string, filters map[string]interface{}) (float64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil ||
		field == "" {
		return 0, errno.InvalidParamError // 如果BaseModel、DB实例或字段为空，返回参数错误
	}

	sField := fmt.Sprintf(`IFNULL(SUM(%s), 0) as total`, field) // 格式化SQL字段
	db = db.Select(sField).Table(m.tableName)                   // 指定查询字段和表名
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val) // 添加过滤条件
		}
	}

	var dao struct{ Total float64 } // 用于存储求和结果
	err := db.Take(&dao).Error      // 执行查询操作
	if OrmErr(err) != nil {
		return 0, err // 如果出错，返回错误
	}
	return dao.Total, nil // 返回求和结果和nil
}

// Distinct 基础模型的去重查询方法
func (m *BaseModel) Distinct(filters map[string]interface{}, fields []string, dao interface{}, page int, pageSize int, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(dao) // 获取dao的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果dao不是指针或为空，返回错误
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order) // 添加排序条件
		}
	}

	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val) // 添加过滤条件
		}
	}
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Distinct(fields).Find(dao).Error // 执行去重查询
	if OrmErr(err) != nil {
		return err // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// Count 基础模型的计数方法
func (m *BaseModel) Count(filters map[string]interface{}) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	db = db.Select("COUNT(1) AS count").Table(m.tableName) // 指定查询字段和表名
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val) // 添加过滤条件
		}
	}

	var dao struct{ Count int64 } // 用于存储计数结果
	err := db.Take(&dao).Error    // 执行查询操作
	if OrmErr(err) != nil {
		return 0, err // 如果出错，返回错误
	}
	return dao.Count, nil // 返回计数结果和nil
}

// ListAndGroupBy 基础模型的分组查询方法
func (m *BaseModel) ListAndGroupBy(filters map[string]interface{}, dao interface{}, fields []string, groupBy string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError // 如果BaseModel或DB实例为空，返回参数错误
	}

	rv := reflect.ValueOf(dao) // 获取dao的反射值
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError // 如果dao不是指针或为空，返回错误
	}

	db = db.Select(fields).Table(m.tableName) // 指定查询字段和表名
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val) // 添加过滤条件
		}
	}
	err := db.Debug().Group(groupBy).Find(dao).Error // 执行分组查询
	if OrmErr(err) != nil {
		return err // 如果出错，返回错误
	}
	return nil // 成功返回nil
}

// DistinctCount 基础模型的去重计数方法
func (m *BaseModel) DistinctCount(filters map[string]interface{}, keys ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db // 如果BaseModel中有DB实例，则使用它
	} else {
		db = NewOrm().Table(m.tableName) // 否则初始化一个新的DB实例并指定表名
	}
	distinctColumn := "" // 用于存储去重字段
	if len(keys) > 0 {
		for _, v := range keys {
			distinctColumn = distinctColumn + v + "," // 拼接去重字段
		}
		distinctColumn = strings.TrimSuffix(distinctColumn, ",") // 去除最后一个逗号
	}
	queryStr := "COUNT(DISTINCT " + distinctColumn + ") AS count" // 格式化SQL查询

	db = db.Select(queryStr).Table(m.tableName).Session(&gorm.Session{}) // 指定查询字段和表名
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val) // 添加过滤条件
		}
	}

	var dao struct{ Count int64 } // 用于存储计数结果
	err := db.Take(&dao).Error    // 执行查询操作
	if OrmErr(err) != nil {
		return constants.DefaultZero, err // 如果出错，返回错误
	}
	return dao.Count, nil // 返回计数结果和nil
}

// OrmErr 处理GORM错误
func OrmErr(err error) error {
	if err != gorm.ErrRecordNotFound {
		return err // 如果不是记录未找到错误，返回错误
	}
	return nil // 否则返回nil
}

// IsUniqueErr 检查是否为唯一约束冲突错误
func IsUniqueErr(err error) bool {
	if err == nil {
		return false // 如果错误为空，返回false
	}
	return strings.Contains(err.Error(), constants.DefaultMysqlUniqueErr) // 检查错误信息是否包含唯一约束冲突
}

// TxBegin 开始事务
func TxBegin(alias ...string) (*gorm.DB, error) {
	tx := NewOrm(alias...).Begin() // 初始化一个新的DB实例并开始事务
	if tx.Error != nil {
		return nil, tx.Error // 如果出错，返回错误
	}
	return tx, nil // 返回事务实例和nil
}

// TxCommit 提交事务
func TxCommit(tx *gorm.DB) error {
	if tx == nil {
		return nil // 如果事务为空，返回nil
	}
	tx.Commit() // 提交事务
	if tx.Error != nil {
		return tx.Error // 如果出错，返回错误
	}
	return nil // 成功返回nil
}
