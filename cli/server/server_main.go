package main

import (
	"errors"
	"os"
	"fmt"

	toml "github.com/BurntSushi/toml"
	server "github.com/z4yx/tunaworks/server"
	cli "github.com/urfave/cli"
)

var gitrev = ""

func loadConfig(cfgFile string, cfg *server.Config) error {
	if cfgFile != "" {
		if _, err := toml.DecodeFile(cfgFile, cfg); err != nil {
			return err
		}
	}

	return nil
}

func parseSettings(c *cli.Context, cfg *server.Config) error {
	if c.Bool("help") {
		cli.ShowAppHelpAndExit(c, 0)
	}
	cf := c.GlobalString("config-file")
	if len(cf) == 0 {
		return errors.New("Config file not specified")
	}
	return loadConfig(cf, cfg)
}

func cmdRun(c *cli.Context) error {
	cfg := server.Config{}
	err := parseSettings(c, &cfg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s := server.MakeServer(&cfg)
	s.Run()
	return nil
}

func main() {
	app := &cli.App{
		Name:      "tunaworks-server",
		UsageText: `tunaworks-server [options]`,
		Usage:     "TUNA.works server",
		Version:   gitrev,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config-file, c", Usage: "`path` to your config file"},
		},
		Action: cmdRun,
		Authors: []cli.Author{
			{Name: "Yuxiang Zhang", Email: "yuxiang.zhang@tuna.tsinghua.edu.cn"},
		},
	}

	app.Run(os.Args)
}
