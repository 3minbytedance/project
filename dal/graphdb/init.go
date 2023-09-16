package graphdb

import (
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const (
	address   = "112.124.58.44"
	port      = 9669
	username  = "root"
	password  = "nebula"
	namespace = "test"
)

// Initialize logger
var log = nebula.DefaultLogger{}
var sessionPool *nebula.SessionPool

func init() {
	hostAddress := nebula.HostAddress{Host: address, Port: port}

	// Create configs for session pool
	config, err := nebula.NewSessionPoolConf(
		username,
		password,
		[]nebula.HostAddress{hostAddress},
		namespace,
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create session pool config, %s", err.Error()))
	}

	// create session pool
	sessionPool, err = nebula.NewSessionPool(*config, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize session pool, %s", err.Error()))
	}
}
