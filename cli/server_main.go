package main

import server "github.com/z4yx/tunaworks/server"

func main() {
	cfg := server.Config{
		Debug: true,
		Server: server.ServerConfig{
			Addr:       "127.0.0.1",
			Port:       8000,
			DBProvider: "sqlite3",
			DBName:     "tunaworks.db",
		},
	}
	s := server.MakeServer(&cfg)
	s.Run()
}
