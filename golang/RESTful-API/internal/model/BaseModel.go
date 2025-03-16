// 封装了DB操作CURD基本方法，用于被其它数据模型继承
package model

import (
	"RESTful-API/internal/constants"
	"RESTful-API/internal/errno"
	"reflect"
	"strings"

	"fmt"

	"gorm.io/gorm"
)

var OrmMap map[string]*gorm.DB

type BaseModel struct {
	db        *gorm.DB
	tableName string
	T         string
}

// NewOrm model new orm
func NewOrm(alias ...string) *gorm.DB {
	ormName := "default"
	if len(alias) > 0 {
		ormName = alias[0]
	}

	if OrmMap == nil {
		panic("[NewOrm] orm not init")
	}

	db, ok := OrmMap[ormName]
	if !ok {
		panic("[NewOrm] orm not exists")
	}

	return db
}

// Create base model insert method
func (m *BaseModel) Create(dao interface{}) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(dao)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}
	db = db.Create(dao)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// Create base model delete method
func (m *BaseModel) Delete(filters map[string]interface{}, dao interface{}) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(dao)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Delete(dao)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// Create base model query one row method
func (m *BaseModel) QueryOne(filters map[string]interface{}, dao interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(dao)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	err := db.Limit(1).Find(dao).Error
	if OrmErr(err) != nil {
		return err
	}
	return nil
}

// Update base model update
func (m *BaseModel) Update(data map[string]interface{}, filters map[string]interface{}, limitNum ...int) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil || db == nil ||
		len(filters) <= 0 ||
		len(data) <= 0 {
		return 0, errno.InvalidParamError
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}

	if len(limitNum) > constants.DefaultZero {
		db = db.Limit(limitNum[0]).Updates(data)
	} else {
		db = db.Updates(data)
	}

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// Rows
func (m *BaseModel) ListAndTotal(filters map[string]interface{}, page, pageSize int, list interface{}, orderBy ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError
	}

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errno.NonPointerOrEmptyError
	}

	var total int64

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Count(&total)
	if OrmErr(db.Error) != nil {
		return 0, db.Error
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(list)
	if OrmErr(db.Error) != nil {
		return 0, db.Error
	}
	return total, nil
}
func (m *BaseModel) TotalByJoin(filters map[string]interface{}, joinSQL string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName).Joins(joinSQL)
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError
	}
	var total int64
	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Count(&total)
	if OrmErr(db.Error) != nil {
		return 0, db.Error
	}

	return total, nil
}

func (m *BaseModel) ListAndTotalByJoin(filters map[string]interface{}, fields []string, joinSQL string, page, pageSize int, list interface{}, orderBy ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError
	}

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errno.NonPointerOrEmptyError
	}

	var total int64

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Count(&total)
	if OrmErr(db.Error) != nil {
		return 0, db.Error
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Select(fields).Joins(joinSQL).Find(list)
	if OrmErr(db.Error) != nil {
		return 0, db.Error
	}
	return total, nil
}

// Rows
func (m *BaseModel) List(filters map[string]interface{}, page, pageSize int, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(list)
	if OrmErr(db.Error) != nil {
		return db.Error
	}
	return nil
}

// Rows
func (m *BaseModel) ListWithJoin(filters map[string]interface{}, fields []string, joinSQL string, page, pageSize int, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize).Select(fields).Joins(joinSQL).Find(list)
	if OrmErr(db.Error) != nil {
		return db.Error
	}
	return nil
}

// list by filter,no page
func (m *BaseModel) ListNoPage(filters map[string]interface{}, list interface{}, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, value := range filters {
			db = db.Where(key, value)
		}
	}

	db = db.Find(list)
	if OrmErr(db.Error) != nil {
		return db.Error
	}
	return nil
}

// Sum
func (m *BaseModel) Sum(field string, filters map[string]interface{}) (float64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil ||
		field == "" {
		return 0, errno.InvalidParamError
	}

	sField := fmt.Sprintf(`IFNULL(SUM(%s), 0) as total`, field)
	db = db.Select(sField).Table(m.tableName)
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val)
		}
	}

	var dao struct{ Total float64 }
	err := db.Take(&dao).Error
	if OrmErr(err) != nil {
		return 0, err
	}
	return dao.Total, nil
}

// distinct
// filter查询条件、fields 查询字段、dao 接收对象
func (m *BaseModel) Distinct(filters map[string]interface{}, fields []string, dao interface{}, page int, pageSize int, orderBy ...string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(dao)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	if len(orderBy) > 0 {
		for _, order := range orderBy {
			db = db.Order(order)
		}
	}

	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val)
		}
	}
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Distinct(fields).Find(dao).Error
	if OrmErr(err) != nil {
		return err
	}
	return nil
}

// count
func (m *BaseModel) Count(filters map[string]interface{}) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return 0, errno.InvalidParamError
	}

	db = db.Select("COUNT(1) AS count").Table(m.tableName)
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val)
		}
	}

	var dao struct{ Count int64 }
	err := db.Take(&dao).Error
	if OrmErr(err) != nil {
		return 0, err
	}
	return dao.Count, nil
}

func (m *BaseModel) ListAndGroupBy(filters map[string]interface{}, dao interface{}, fields []string, groupBy string) error {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	if m == nil ||
		db == nil {
		return errno.InvalidParamError
	}

	rv := reflect.ValueOf(dao)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errno.NonPointerOrEmptyError
	}

	db = db.Select(fields).Table(m.tableName)
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val)
		}
	}
	err := db.Debug().Group(groupBy).Find(dao).Error
	if OrmErr(err) != nil {
		return err
	}
	return nil
}

// distinct count
func (m *BaseModel) DistinctCount(filters map[string]interface{}, keys ...string) (int64, error) {
	var db *gorm.DB
	if m.db != nil {
		db = m.db
	} else {
		db = NewOrm().Table(m.tableName)
	}
	distinctColumn := ""
	if len(keys) > 0 {
		for _, v := range keys {
			distinctColumn = distinctColumn + v + ","
		}
		distinctColumn = strings.TrimSuffix(distinctColumn, ",")
	}
	queryStr := "COUNT(DISTINCT " + distinctColumn + ") AS count"

	db = db.Select(queryStr).Table(m.tableName).Session(&gorm.Session{})
	if len(filters) > 0 {
		for key, val := range filters {
			db = db.Where(key, val)
		}
	}

	var dao struct{ Count int64 }
	err := db.Take(&dao).Error
	if OrmErr(err) != nil {
		return constants.DefaultZero, err
	}
	return dao.Count, nil
}

// OrmErr return orm error
func OrmErr(err error) error {
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

// IsUniqueErr check if unique conflict
func IsUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), constants.DefaultMysqlUniqueErr)
}

// TxBegin common tx begingithub.com/go-sql-driver/mysql@latest
func TxBegin(alias ...string) (*gorm.DB, error) {
	tx := NewOrm(alias...).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// TxCommit common tx commit
func TxCommit(tx *gorm.DB) error {
	if tx == nil {
		return nil
	}
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
