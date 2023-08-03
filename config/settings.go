package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Conf      = new(AppConfig)
	LocalMode = "remote"
)

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`

	Local struct {
		*MySQLConfig `mapstructure:"mysql"`
		*RedisConfig `mapstructure:"redis"`
	} `mapstructure:"local"`

	Remote struct {
		*MySQLConfig `mapstructure:"mysql"`
		*RedisConfig `mapstructure:"redis"`
	} `mapstructure:"remote"`
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
	Address      string `mapstructure:"address"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	VCIdDB       int    `mapstructure:"vcid_db"`
	CVIdDB       int    `mapstructure:"cvid_db"`
	CIdCommentDB int    `mapstructure:"cid_comment_db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func Init() (err error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("app")
	//viper.SetConfigFile("config.yaml") // 指定配置文件路径
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {            // 读取配置信息失败
		panic(fmt.Errorf("Read config.yaml failed: %s \n", err))
	}

	// 读取到的配置信息 反序列化到 Conf 里面
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("Viper unmarshal failed: %v\n", err)
	}

	// 监控配置文件变化, 实时更新Conf
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变化了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("Viper unmarshal failed, err: %v\n", err)
		}
	})
	return
}