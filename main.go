package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	ydns "github.com/wyattjoh/ydns-updater/internal"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "ydns-updater"
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "base",
			Value: "https://ydns.io/api/v1/update/",
			Usage: "base url for api calls on ydns",
		},
		&cli.StringFlag{
			Name:     "host",
			EnvVars:  []string{"YDNS_HOST"},
			Required: true,
			Usage:    "host to update",
		},
		&cli.StringFlag{
			Name:     "user",
			EnvVars:  []string{"YDNS_USER"},
			Required: true,
			Usage:    "username for authentication on ydns",
		},
		&cli.StringFlag{
			Name:     "pass",
			EnvVars:  []string{"YDNS_PASS"},
			Required: true,
			Usage:    "password for authentication on ydns",
		},
		&cli.BoolFlag{
			Name:    "daemon",
			EnvVars: []string{"YDNS_DAEMON"},
			Usage:   "enables the updater as a daemon",
		},
		&cli.DurationFlag{
			Name:    "frequency",
			EnvVars: []string{"YDNS_FREQUENCY"},
			Value:   60 * time.Minute,
			Usage:   "sleep time between updates while in daemon mode",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enables debug logging",
		},
	}
	app.Action = func(c *cli.Context) error {
		base := c.String("base")
		host := c.String("host")
		user := c.String("user")
		pass := c.String("pass")
		daemon := c.Bool("daemon")
		frequency := c.Duration("frequency")

		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if err := ydns.Run(base, host, user, pass); err != nil {
			return cli.Exit(err, 1)
		}

		for daemon {
			logrus.WithField("sleep", frequency).Info("sleeping till next update")
			time.Sleep(frequency)

			if err := ydns.Run(base, host, user, pass); err != nil {
				return cli.Exit(err, 1)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal()
	}
}
