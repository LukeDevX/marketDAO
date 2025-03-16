package bootstrap

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"marketDAO/golang/RESTful-API/utils/config"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		initDB()
	})
}

func initDB() {
	if enable, _ := config.GetConfig("db.enable").Bool(); enable == false {
		return
	}

	logs.Info("Start initial DB...")

	model.OrmMap = make(map[string]*gorm.DB)
	dbAlias := config.GetConfig("db.alias").Strings(",")
	if len(dbAlias) <= constants.DefaultZero {
		panic("Invalid db config")
	}

	for _, alias := range dbAlias {
		host := config.GetConfig("db." + alias + ".host").String()
		port := config.GetConfig("db." + alias + ".port").String()
		user := config.GetConfig("db." + alias + ".user").String()
		pwd := config.GetConfig("db." + alias + ".password").String()
		dbName := config.GetConfig("db." + alias + ".name").String()
		timezone := config.GetConfig("db." + alias + ".timezone").String()

		maxIdle, _ := config.GetConfig("db." + alias + ".maxIdle").Int()
		maxConn, _ := config.GetConfig("db." + alias + ".maxConnections").Int()
		maxLifeTime, _ := config.GetConfig("db." + alias + ".maxLifeTime").Int()
		maxLifeDuration := time.Second * time.Duration(maxLifeTime)

		debug, _ := config.GetConfig("db." + alias + ".debug").Bool()

		// default port
		if port == "" {
			port = "3306"
		}
		// dns
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, pwd, host, port, dbName)
		if timezone != "" {
			dsn = dsn + "&loc=" + url.QueryEscape(timezone)
		}

		var newLogger logger.Interface
		//调试模式
		if debug {
			newLogger = logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
					LogLevel:      logger.Info,            // Log level
					Colorful:      true,                   // 彩色打印
				},
			)
		}
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		sqlDB.SetConnMaxLifetime(maxLifeDuration)
		sqlDB.SetMaxOpenConns(maxConn)
		sqlDB.SetMaxIdleConns(maxIdle)
		model.OrmMap[alias] = db

		logs.Info("DB Finished Start %v", dsn)
	}
}
