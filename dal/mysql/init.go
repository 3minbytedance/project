package mysql

import (
	"douyin/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.MySQLConfig
	if appConfig.Mode == config.LocalMode {
		conf = appConfig.Local.MySQLConfig
	} else {
		conf = appConfig.Remote.MySQLConfig
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Address,
		conf.Port,
		conf.Database,
	)

	mysqlLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: mysqlLog})
	if err != nil {
		fmt.Println(dsn)
		log.Fatal("connect to mysql failed:", err)
	}
	//err = DB.AutoMigrate(&models.User{})
	//if err != nil {
	//	return
	//}
	return nil
}
