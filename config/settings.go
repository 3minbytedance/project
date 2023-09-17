package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var (
	Conf      = new(AppConfig)
	LocalMode = "local"
)

type AppConfig struct {
	Name       string `mapstructure:"name"`
	Node       string `mapstructure:"node"`
	Mode       string `mapstructure:"mode"`
	Port       int    `mapstructure:"port"`
	Version    string `mapstructure:"version"`
	*LogConfig `mapstructure:"log"`

	Local struct {
		*MySQLConfig    `mapstructure:"mysql"`
		*RedisConfig    `mapstructure:"redis"`
		*KafkaConfig    `mapstructure:"kafka"`
		*MongoConfig    `mapstructure:"mongo"`
		*RocketMQConfig `mapstructure:"rocketmq"`
		*GraphDBConfig  `mapstructure:"graphDB"`
	} `mapstructure:"local"`

	Remote struct {
		*MySQLConfig    `mapstructure:"mysql"`
		*RedisConfig    `mapstructure:"redis"`
		*KafkaConfig    `mapstructure:"kafka"`
		*MongoConfig    `mapstructure:"mongo"`
		*RocketMQConfig `mapstructure:"rocketmq"`
		*GraphDBConfig  `mapstructure:"graphDB"`
	} `mapstructure:"remote"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type RedisConfig struct {
	Address           string `mapstructure:"address"`
	Port              int    `mapstructure:"port"`
	Password          string `mapstructure:"password"`
	DB                int    `mapstructure:"db"`
	UerFavoriteRDB    int    `mapstructure:"ufvdb"`
	VideoFavoritedRDB int    `mapstructure:"vfudb"`
	PoolSize          int    `mapstructure:"pool_size"`
	MinIdleConns      int    `mapstructure:"min_idle_conns"`
	CommentDB         int    `mapstructure:"comment_db"`
	ExpireTime        int64  `mapstructure:"expire_time"`
}

type KafkaConfig struct {
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MongoConfig struct {
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type RocketMQConfig struct {
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type GraphDBConfig struct {
	Address   string `mapstructure:"address"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Namespace string `mapstructure:"namespace"`
}

func Init() (err error) {
	//viper.AddConfigPath("../../config")
	//viper.SetConfigName("app")
	viper.SetConfigFile("config/app.yaml") // 指定配置文件路径
	err = viper.ReadInConfig()             // 读取配置信息
	if err != nil {                        // 读取配置信息失败¬
		log.Fatalf("Read app.yaml failed: %s \n", err)
	}

	// 读取到的配置信息 反序列化到 Conf 里面
	if err := viper.Unmarshal(Conf); err != nil {
		log.Printf("Viper unmarshal failed: %v\n", err)
	}

	// 监控配置文件变化, 实时更新Conf
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置发生变化了...")
		if err := viper.Unmarshal(Conf); err != nil {
			log.Printf("Viper unmarshal failed, err: %v\n", err)
		}
	})

	return
}
