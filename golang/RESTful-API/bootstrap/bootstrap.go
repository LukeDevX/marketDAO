// Package bootstrap 用于初始化系统启动时的必要组件
package bootstrap

// 引入项目内部及外部的必要依赖
import (
	"RESTful-API/internal/constants" // 常量定义
	"RESTful-API/internal/model"     // 数据库 ORM 模型
	"RESTful-API/utils/config"       // 配置管理
	"RESTful-API/utils/logs"         // 日志工具
	"fmt"                            // 标准库: 格式化输入输出
	"gorm.io/driver/mysql"           // GORM 的 MySQL 适配器
	"gorm.io/gorm"                   // GORM ORM 框架
	"gorm.io/gorm/logger"            // GORM 日志
	"log"                            // 标准库: 日志
	"net/url"                        // 标准库: 处理 URL 相关操作
	"os"                             // 标准库: 操作系统相关
	"sync"                           // 标准库: 并发同步
	"time"                           // 标准库: 时间处理
)

// 定义 `once` 变量，确保 `init()` 只执行一次（单例模式）
var (
	once sync.Once
)

// init() 函数，在包初始化时执行，确保 `initDB()` 只执行一次
func init() {
	once.Do(func() { // `sync.Once` 确保 `initDB()` 仅运行一次
		initLog()
		initDB()
	})
}

// 初始化日志
func initLog() {
	logs.InitLog()
}

// 初始化数据库连接
func initDB() {
	// 从配置文件读取数据库启用状态，如果 `db.enable` 为 false，则直接返回
	if enable, _ := config.GetConfig("db.enable").Bool(); enable == false {
		return
	}

	// 记录日志，表示开始初始化数据库
	logs.Info("Start initial DB...")
	fmt.Println("Start initial DB...")

	// 创建数据库连接的全局映射，存储多个数据库连接
	model.OrmMap = make(map[string]*gorm.DB)

	// 获取数据库别名列表，例如 ["default", "backup"]
	dbAlias := config.GetConfig("db.alias").Strings(",")

	// 检查是否至少有一个数据库配置，否则抛出异常
	if len(dbAlias) <= constants.DefaultZero {
		panic("Invalid db config")
	}

	// 遍历数据库别名列表，逐个连接数据库
	for _, alias := range dbAlias {
		// 从配置文件读取数据库连接参数
		host := config.GetConfig("db." + alias + ".host").String()
		port := config.GetConfig("db." + alias + ".port").String()
		user := config.GetConfig("db." + alias + ".user").String()
		pwd := config.GetConfig("db." + alias + ".password").String()
		dbName := config.GetConfig("db." + alias + ".name").String()
		timezone := config.GetConfig("db." + alias + ".timezone").String()

		// 读取数据库连接池参数
		maxIdle, _ := config.GetConfig("db." + alias + ".maxIdle").Int()         // 最大空闲连接数
		maxConn, _ := config.GetConfig("db." + alias + ".maxConnections").Int()  // 最大连接数
		maxLifeTime, _ := config.GetConfig("db." + alias + ".maxLifeTime").Int() // 连接最大存活时间
		maxLifeDuration := time.Second * time.Duration(maxLifeTime)              // 转换为 `time.Duration`

		// 是否开启 GORM 调试模式（输出 SQL 语句）
		debug, _ := config.GetConfig("db." + alias + ".debug").Bool()

		// 如果端口为空，使用 MySQL 默认端口 3306
		if port == "" {
			port = "3306"
		}

		// 组装 MySQL 连接字符串（DSN: Data Source Name）
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, pwd, host, port, dbName)
		// 如果配置了时区，则添加到 DSN
		if timezone != "" {
			dsn = dsn + "&loc=" + url.QueryEscape(timezone)
		}

		// 定义 GORM 日志对象
		var newLogger logger.Interface
		// 如果启用 debug 模式，则启用 SQL 语句日志
		if debug {
			newLogger = logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出到标准输出（控制台）
				logger.Config{
					SlowThreshold: 200 * time.Millisecond, // 慢查询阈值（200ms）
					LogLevel:      logger.Info,            // 日志级别
					Colorful:      true,                   // 彩色日志
				},
			)
		}

		// 连接数据库，使用 GORM + MySQL 驱动
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			panic(err) // 连接失败则直接退出
		}

		// 获取底层 SQL 连接
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		// 配置数据库连接池参数
		sqlDB.SetConnMaxLifetime(maxLifeDuration) // 连接最大存活时间
		sqlDB.SetMaxOpenConns(maxConn)            // 最大连接数
		sqlDB.SetMaxIdleConns(maxIdle)            // 最大空闲连接数

		// 将数据库连接存入全局映射 `OrmMap`
		model.OrmMap[alias] = db

		// 记录数据库启动完成日志
		logs.Info("DB Finished Start %v", dsn)
		fmt.Println("Start dsn", dsn)
	}
}
