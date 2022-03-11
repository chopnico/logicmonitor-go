package main

import (
	"context"
	"errors"
	"log"
	"os"

	logicmonitor "github.com/chopnico/logicmonitor-go"

	CLI "github.com/chopnico/logicmonitor-go/internal/cli"

	"github.com/urfave/cli/v2"
)

// some application and default variables
var (
	AppName  string = "logicmonitor"
	AppUsage string = "a logicmonitor cli/tui tool"
	// ldflags will be used to set this. check Makefile
	AppVersion string

	DefaultLoggingLevel = "info"
	DefaultPrintFormat  = "table"
	DefaultTimeOut      = 60
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "accessID",
			Usage:       "account `ACCESSID`",
			EnvVars:     []string{"LOGICMONITOR_ACCESSID"},
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "accessKey",
			Usage:       "account `ACCESSKEY`",
			EnvVars:     []string{"LOGICMONITOR_ACCESSKEY"},
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "company",
			Usage:       "logicmonitor appliance `COMPANY`",
			EnvVars:     []string{"LOGICMONITOR_COMPANY"},
			DefaultText: "none",
		},
		&cli.BoolFlag{
			Name:    "ignore-ssl",
			Usage:   "ignore ssl errors",
			EnvVars: []string{"LOGICMONITOR_IGNORE_SSL"},
			Value:   false,
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "http timeout",
			Value: 0,
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "printing format (json, list, table)",
			Value: "table",
		},
		&cli.StringFlag{
			Name:  "logging",
			Usage: "set logging level",
			Value: "info",
		},
		&cli.StringFlag{
			Name:    "proxy",
			Usage:   "set http proxy",
			EnvVars: []string{"LOGICMONITOR_PROXY"},
		},
	}
	app.Before = func(c *cli.Context) error {
		var err error
		var api *logicmonitor.API

		if c.String("accessID") == "" {
			cli.ShowAppHelp(c)
			return errors.New(logicmonitor.ErrorEmptyAccessID)
		} else if c.String("accessKey") == "" {
			cli.ShowAppHelp(c)
			return errors.New(logicmonitor.ErrorEmptyAccessKey)
		} else if c.String("company") == "" {
			cli.ShowAppHelp(c)
			return errors.New(logicmonitor.ErrorEmptyCompany)
		}

		api, err = logicmonitor.NewAPIBasicAuth(
			c.String("accessID"),
			c.String("accessKey"),
			c.String("company"),
		)
		if err != nil {
			return err
		}

		// set options
		api.Timeout(c.Int("timeout")).
			LoggingLevel(c.String("logging")).
			Proxy(c.String("proxy"))

		// should we ignore ssl errors?
		if c.Bool("ignore-ssl") {
			api.IgnoreSSLErrors()
		}

		ctx := context.WithValue(c.Context, logicmonitor.APIContextKey("api"), api)
		c.Context = ctx

		return nil
	}

	// create cli commands
	CLI.NewCommands(app)

	// run the app
	err := app.Run(os.Args)
	if err != nil {
		if err.Error() != "debugging" {
			log.Fatal(err)
		}
	}
	os.Exit(0)
}
