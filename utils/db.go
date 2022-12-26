package utils

import (
	"sync"

	"go-pano/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IDBInstance interface {
	DB() *gorm.DB
}

// DBInstance is a singleton DB instance
type DBInstance struct {
	//
	initializer func() interface{}
	instance    interface{}
	once        sync.Once
}

var (
	dbInstance *DBInstance
)

// 獲取實例，且避免重複實例化
func (i *DBInstance) Instance() interface{} {
	i.once.Do(func() {
		i.instance = i.initializer()
	})
	return i.instance
}

// DB初始化
func dbInit() interface{} {
	lv := logger.Error
	if config.Server.Mode != gin.ReleaseMode {
		lv = logger.Info // output debug logs in dev mode
	}

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(lv),
	}

	db, err := gorm.Open(mysql.Open(config.Database.DSN), cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to database")
	}

	stdDB, _ := db.DB()
	stdDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	stdDB.SetMaxOpenConns(config.Database.MaxOpenConns)

	return db
}

// 連線成功後獲取DB的連線訊息
func DB() *gorm.DB {
	return dbInstance.Instance().(*gorm.DB)
}

// 初始化
func init() {
	dbInstance = &DBInstance{initializer: dbInit}
}
