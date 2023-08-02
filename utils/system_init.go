package utils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DB    *gorm.DB
	Red   *redis.Client
	local = true
)

func InitConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("app")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read configs failed!")
	}
}

func InitMysql() {
	mysqlLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	address := os.Getenv("MYSQL_HOST")
	portStr := os.Getenv("MYSQL_PORT")
	port, _ := strconv.Atoi(portStr)
	fmt.Println("address: ", address, " port: ", port)
	dsn := ""
	if local {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds",
			viper.GetString("local.mysql.username"),
			viper.GetString("local.mysql.password"),
			viper.GetString("local.mysql.address"),
			viper.GetInt("local.mysql.port"),
			viper.GetString("local.mysql.database"),
			viper.GetInt("local.mysql.timeout"))
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds",
			viper.GetString("remote.mysql.username"),
			viper.GetString("remote.mysql.password"),
			//viper.GetString("remote.mysql.address"),
			//viper.GetInt("remote.mysql.port"),
			address,
			port,
			viper.GetString("remote.mysql.database"),
			viper.GetInt("remote.mysql.timeout"))
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: mysqlLog})
	if err != nil {
		fmt.Println(dsn)
		log.Fatal("connect to mysql failed")
	}
	DB = db
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("local.redis.addr"),
		Password:     viper.GetString("local.redis.password"),
		DB:           viper.GetInt("local.redis.DB"),
		PoolSize:     viper.GetInt("local.redis.poolSize"),
		MinIdleConns: viper.GetInt("local.redis.minIdleConn"),
	})
}
