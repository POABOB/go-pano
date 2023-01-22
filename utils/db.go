package utils

import (
	"fmt"
	"sync"

	"go-pano/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *DBInstance
)

// func init() {
// 	dbInstance = &DBInstance{instance: dbInit()}
// }

type IDBInstance interface {
	DB() *gorm.DB
}

// DBInstance is a singleton DB instance
type DBInstance struct {
	instance *gorm.DB
	once     sync.Once
}

// 獲取實例，且避免重複實例化
func (i *DBInstance) DB() *gorm.DB {
	i.once.Do(func() {
		fmt.Println(12356)
		i.instance = dbInit()
	})
	return i.instance
}

// 獲取實例，且避免重複實例化
func NewMockInstance(m *gorm.DB, o sync.Once) IDBInstance {
	o.Do(func() {})
	return &DBInstance{m, o}
}

// DB初始化
func dbInit() *gorm.DB {
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
