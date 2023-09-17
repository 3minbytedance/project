package graphdb

import (
	"douyin/config"
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

// Initialize logger
var log = nebula.DefaultLogger{}
var sessionPool *nebula.SessionPool

func Init(appConfig *config.AppConfig) (err error) {
	var conf *config.GraphDBConfig
	if appConfig.Mode == config.LocalMode {
		conf = appConfig.Local.GraphDBConfig
	} else {
		conf = appConfig.Remote.GraphDBConfig
	}

	hostAddress := nebula.HostAddress{Host: conf.Address, Port: conf.Port}

	// Create configs for session pool
	configs, err := nebula.NewSessionPoolConf(
		conf.Username,
		conf.Password,
		[]nebula.HostAddress{hostAddress},
		conf.Namespace,
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create graphDB session pool config, %s", err.Error()))
		return err
	}

	// create session pool
	sessionPool, err = nebula.NewSessionPool(*configs, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize graphDB session pool, %s", err.Error()))
		return err
	}

	return nil
}
