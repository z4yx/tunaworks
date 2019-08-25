package server

import (
	_ "github.com/BurntSushi/toml"
)

// A Config is the top-level toml-serializaible config struct
type Config struct {
	Debug  bool         `toml:"debug"`
	Server ServerConfig `toml:"server"`
}

// A ServerConfig represents the configuration for HTTP server
type ServerConfig struct {
	Addr       string `toml:"addr"`
	Port       int    `toml:"port"`
	DBProvider string `toml:"db_provider"`
	DBName     string `toml:"db_name"`
}
