package prober

import (
	_ "github.com/BurntSushi/toml"
)

type ProberConfig struct {
	Debug    bool   `toml:"debug"`
	IPv4     bool   `toml:"ipv4"`
	IPv6     bool   `toml:"ipv6"`
	Server   string `toml:"server"`
	Https    bool   `toml:"https"`
	Token    string `toml:"token"`
	Interval int    `toml:"interval"`
	Version  string
}
