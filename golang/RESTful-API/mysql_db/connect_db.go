package mysql_db

// package main // 这个注释掉的行表示该文件原本可以作为独立的 Go 可执行程序，但目前作为包导入使用

import (
	"fmt"                    // 标准库，提供格式化 I/O 功能，例如打印日志和调试输出
	"github.com/spf13/viper" // Viper 库，用于读取和解析配置文件（支持 YAML、JSON、TOML 等）
	"gorm.io/driver/mysql"   // GORM 的 MySQL 连接驱动，用于连接 MySQL 数据库
	"gorm.io/gorm"           // GORM（Go 的 ORM 框架），用于数据库操作
	"gorm.io/gorm/schema"    // GORM 的 Schema 相关功能，比如自定义表名规则
	"time"                   // 标准库，提供时间和日期处理
)

// 连接数据库函数，返回 *gorm.DB 连接实例和可能的错误
func ConnectToDatabase() (*gorm.DB, error) {
	// 连接数据库
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情

	// 构造 MySQL DSN（数据源名称），格式为：
	// "用户名:密码@tcp(主机:端口)/数据库名?charset=编码&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("db.UserName"), // 从 Viper 配置文件获取数据库用户名
		viper.GetString("db.Password"), // 获取数据库密码
		viper.GetString("db.Host"),     // 获取数据库主机地址
		viper.GetString("db.Port"),     // 获取数据库端口
		viper.GetString("db.Database"), // 获取数据库名称
		viper.GetString("db.Charset"))  // 获取数据库字符集

	// 使用 GORM 连接 MySQL 数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置表名为单数形式，例如 `user` 而不是 `users`
		},
	})
	if err != nil {
		return nil, err // 连接失败，返回错误
	}

	// 获取底层 SQL 数据库对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err // 获取失败，返回错误
	}

	// 设置数据库连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池的最大连接数
	sqlDB.SetMaxOpenConns(100)          // 设置最大打开的连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接的最大生命周期，避免连接长期占用资源

	return db, nil // 返回数据库连接实例
}

// 数据库自动迁移函数
func CreateDB(db *gorm.DB, model *any) error {
	// 使用 GORM 的 AutoMigrate 进行数据库迁移（自动创建表结构）
	err := db.AutoMigrate(model)
	return err
}

// init 函数：Go 语言的特殊函数，在包初始化时自动执行
func init() {
	// 这里可以放置初始化逻辑，例如加载配置文件
}

// 主函数
func main() {
	// 连接数据库
	db, err := ConnectToDatabase()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	// 数据库迁移（创建 `CrudList` 表）
	err = db.AutoMigrate(&CrudList{})
	if err != nil {
		fmt.Println("Error migrating database:", err)
	}
}
