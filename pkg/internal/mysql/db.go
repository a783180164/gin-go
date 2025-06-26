package mysql

import (
	"fmt"
	"sync"
	"time"

	"gin-go/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	once sync.Once
	db   *gorm.DB
	err  error
)

type Predicate string

var (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
	InPredicate                 = Predicate("IN")
)

var config = configs.Get()

// Init 初始化数据库连接，建议在 main() 一启动就调用
func Init() error {
	once.Do(func() {
		cfg := config.MySQL

		// 构造 write DSN
		writeDSN := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Write.User, cfg.Write.Pass,
			cfg.Write.Addr, cfg.Write.Name,
		)

		// 构造 read DSN
		readDSN := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Read.User, cfg.Read.Pass,
			cfg.Read.Addr, cfg.Read.Name,
		)

		// 打开主库（write）
		db, err = gorm.Open(mysql.Open(writeDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return
		}

		// 装载读写分离插件，Replica 为 read 库
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(writeDSN)},
			Replicas: []gorm.Dialector{mysql.Open(readDSN)},
			// 默认策略：随机选择 Replica
		}))
		if err != nil {
			return
		}

		// 设置连接池参数
		sqlDB, err2 := db.DB()
		if err2 != nil {
			err = err2
			return
		}
		sqlDB.SetMaxOpenConns(cfg.Base.MaxOpenConn)
		sqlDB.SetMaxIdleConns(cfg.Base.MaxIdleConn)
		sqlDB.SetConnMaxLifetime(cfg.Base.ConnMaxLifeTime * time.Second)
	})

	return err
}

// Instance 返回全局 *gorm.DB
// 请确保在调用 Instance 之前已经运行过 Init()
func Instance() *gorm.DB {
	if db == nil {
		panic("database not initialized: call db.Init() first")
	}
	return db
}
