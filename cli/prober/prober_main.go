package main

import (
	"errors"
	"os"
	"fmt"

	toml "github.com/BurntSushi/toml"
	prober "github.com/z4yx/tunaworks/prober"
	"gopkg.in/op/go-logging.v1"
	cli "gopkg.in/urfave/cli.v1"
)

var logger = logging.MustGetLogger("tunaworks")

func loadConfig(cfgFile string, cfg *prober.ProberConfig) error {
	if cfgFile != "" {
		if _, err := toml.DecodeFile(cfgFile, cfg); err != nil {
			return err
		}
	}

	return nil
}

func parseSettings(c *cli.Context, cfg *prober.ProberConfig) error {
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
	cfg := prober.ProberConfig{}
	err := parseSettings(c, &cfg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s := prober.MakeProber(&cfg)
	s.Run()
	return nil
}

func main() {
	app := &cli.App{
		Name:      "tunaworks-prober",
		UsageText: `tunaworks-prober [options]`,
		Usage:     "TUNA.works prober",
		Version:   "1.0",
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
