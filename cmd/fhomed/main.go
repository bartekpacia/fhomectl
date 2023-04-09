package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var k = koanf.New(".")

func main() {
	app := &cli.App{
		Name:  "fhomed",
		Usage: "Long-running daemon for F&Home Cloud",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output logs in JSON Lines format",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "show debug logs",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of the HomeKit bridge accessory",
				Value: "fhomed",
			},
			&cli.StringFlag{
				Name:  "pin",
				Usage: "PIN of the HomeKit bridge accessory",
				Value: "00102003",
			},
		},
		Before: before,
		Action: func(c *cli.Context) error {
			name := c.String("name")
			pin := c.String("pin")

			return daemon(name, pin)
		},
		CommandNotFound: func(c *cli.Context, command string) {
			log.Printf("invalid command '%s'. See 'fhomed --help'\n", command)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("exit", slog.Any("error", err))
		os.Exit(1)
	}
}

func before(c *cli.Context) error {
	var level slog.Level
	if c.Bool("debug") {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	if c.Bool("jsonl") {
		logger := slog.New(slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stdout))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(tint.Options{Level: level, TimeFormat: time.TimeOnly}.NewHandler(os.Stdout))
		slog.SetDefault(logger)
	}

	k := koanf.New(".")
	p := "/etc/fhomed/config.toml"
	if err := k.Load(file.Provider(p), toml.Parser()); err != nil {
		slog.Debug("failed to load config file", slog.Any("error", err))
	} else {
		slog.Debug("loaded config file", slog.String("path", p))
	}

	homeDir, _ := os.UserHomeDir()
	p = fmt.Sprintf("%s/.config/fhomed/config.toml", homeDir)
	if err := k.Load(file.Provider(p), toml.Parser()); err != nil {
		slog.Debug("failed to load config file", slog.Any("error", err))
	} else {
		slog.Debug("loaded config file", slog.String("path", p))
	}

	return nil
}
