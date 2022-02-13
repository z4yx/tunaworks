package server

import (
	"gopkg.in/op/go-logging.v1"
)

var logTag = "tunaworks"
var logger = logging.MustGetLogger(logTag)

func setLogLevel(cfg *Config) {
	if cfg.Debug {
		logging.SetLevel(logging.DEBUG, logTag)
	} else {
		logging.SetLevel(logging.WARNING, logTag)
	}
}
