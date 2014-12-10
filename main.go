package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	BaseURL  string
	Host     string
	User     string
	Password string
}

func readConfig(c *cli.Context) (err error, config *Config) {
	config = &Config{}

	// Grab flag
	config.BaseURL = c.String("base")
	if config.BaseURL == "" {
		// Log out error
		err = errors.New("--base not defined.")
		return
	}

	// Grab flag
	config.Host = c.String("host")
	if config.Host == "" {
		// Log out error
		err = errors.New("--host not defined.")
		return
	}

	// Grab flag
	config.User = c.String("user")
	if config.User == "" {
		// Log out error
		err = errors.New("--user not defined.")
		return
	}

	// Grab flag
	config.Password = c.String("pass")
	if config.Password == "" {
		// Log out error
		err = errors.New("--pass not defined.")
		return
	}

	return
}

func performUpdate(c *Config) {
	// Build the url
	updateUrl, err := url.Parse(c.BaseURL)
	if err != nil {
		log.Fatal("Cannot create url:", err)
	}

	// Build the update params
	var updateParams url.Values = make(url.Values)

	// Add the host param
	updateParams.Add("host", c.Host)

	// Add query params to url
	updateUrl.RawQuery = updateParams.Encode()

	// Build a request
	req, err := http.NewRequest("GET", updateUrl.String(), nil)
	if err != nil {
		log.Fatal("Cannot create request:", err)
	}

	// Set authentication
	req.SetBasicAuth(c.User, c.Password)

	// Get the http client
	client := &http.Client{}

	log.Printf("Updating record %s...", c.Host)

	// Ask the client to perform request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Cannot perform http get:", err)
	}
	defer resp.Body.Close()

	// Log based on request status code
	switch resp.StatusCode {
	case 200:
		log.Printf("Update of %s was successful.\n", c.Host)
	case 400:
		log.Fatal("Failed to perform request due to invalid input parameters.")
	case 401:
		log.Fatal("Failed to perform request due to authentication issues.")
	case 404:
		log.Fatal("Failed to perform request because the host you'd like to update cannot be found.")
	default:
		log.Fatal("Some unknown error occured:", resp.Status)
	}

}

func main() {
	app := cli.NewApp()
	app.Name = "ydns-updater"
	app.Usage = "updates dns entries on ydns"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "base",
			Value: "https://ydns.eu/api/v1/update/",
			Usage: "Base url for api calls on ydns.",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "",
			Usage: "Host to update.",
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "API Username for authentication on ynds.",
		},
		cli.StringFlag{
			Name:  "pass",
			Value: "",
			Usage: "API Password for authentication on ynds.",
		},
	}
	app.Action = func(c *cli.Context) {
		// Read config from context
		err, config := readConfig(c)
		if err != nil {
			// Print log
			log.Fatal(err)

			// Exit
			return
		}

		// Perform update
		performUpdate(config)
	}

	app.Run(os.Args)
}
