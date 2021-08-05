package main

import (
	"errors"
	"os"
	"fmt"

	toml "github.com/BurntSushi/toml"
	prober "github.com/z4yx/tunaworks/prober"
	cli "github.com/urfave/cli"
)

var gitrev = ""

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
	cfg.Version = gitrev
	s := prober.MakeProber(&cfg)
	s.Run()
	return nil
}

func main() {
	app := &cli.App{
		Name:      "tunaworks-prober",
		UsageText: `tunaworks-prober [options]`,
		Usage:     "TUNA.works prober",
		Version:   gitrev,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "config-file, c",
				Value: "prober.yaml",
				Usage: "`path` to your config file",
			},
		},
		Action: cmdRun,
		Authors: []cli.Author{
			{Name: "Yuxiang Zhang", Email: "yuxiang.zhang@tuna.tsinghua.edu.cn"},
		},
	}

	app.Run(os.Args)
}
