package main

import server "github.com/z4yx/tunaworks/server"

func main() {
	cfg := server.Config{
		Debug: true,
		Server: server.ServerConfig{
			Addr:       "0.0.0.0",
			Port:       8001,
			DBProvider: "sqlite3",
			DBName:     "tunaworks.db",
		},
	}
	s := server.MakeServer(&cfg)
	s.Run()
}
