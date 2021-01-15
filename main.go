package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func run(base, host, user, pass string) error {
	u, err := url.Parse(base)
	if err != nil {
		return errors.Wrap(err, "cannot create url")
	}

	u.Query().Add("host", host)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	req.SetBasicAuth(user, pass)

	logrus.WithField("host", host).Info("updating record")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot perform http get")
	}
	defer res.Body.Close()

	// Log based on request status code
	switch res.StatusCode {
	case http.StatusOK:
		logrus.WithField("host", host).Info("update was successful")
	case http.StatusBadRequest:
		return errors.New("failed to perform request due to invalid input parameters")
	case http.StatusUnauthorized:
		return errors.New("failed to perform request due to authentication issues")
	case http.StatusNotFound:
		return errors.New("failed to perform request because the host you'd like to update cannot be found")
	default:
		return errors.Errorf("some unknown error occurred: %s", res.Status)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ydns-updater"
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "base",
			Value:    "https://ydns.io/api/v1/update/",
			Required: true,
			Usage:    "base url for api calls on ydns",
		},
		&cli.StringFlag{
			Name:     "host",
			Required: true,
			Usage:    "host to update",
		},
		&cli.StringFlag{
			Name:     "user",
			Required: true,
			Usage:    "username for authentication on ydns",
		},
		&cli.StringFlag{
			Name:     "pass",
			Required: true,
			Usage:    "password for authentication on ydns",
		},
		&cli.BoolFlag{
			Name:  "daemon",
			Usage: "enables the updater as a daemon",
		},
		&cli.IntFlag{
			Name:  "frequency",
			Value: 60,
			Usage: " minutes between updates while in daemon mode",
		},
	}
	app.Action = func(c *cli.Context) error {
		base := c.String("base")
		host := c.String("host")
		user := c.String("user")
		pass := c.String("pass")
		daemon := c.Bool("daemon")
		frequency := c.Int("frequency")

		sleep := time.Duration(frequency) * time.Minute

		if err := run(base, host, user, pass); err != nil {
			return cli.Exit(err, 1)
		}

		for daemon {
			logrus.WithField("sleep", sleep).Info("sleeping till next update")
			time.Sleep(sleep)

			if err := run(base, host, user, pass); err != nil {
				return cli.Exit(err, 1)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal()
	}
}
