package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Config describes the configuration options for the update operation.
type Config struct {
	BaseURL  string
	Host     string
	User     string
	Password string
}

func verifyConfig(config Config) error {
	if config.BaseURL == "" {
		return errors.New("--base not defined, see usage")
	}

	if config.Host == "" {
		return errors.New("--host not defined, see usage")
	}

	if config.User == "" {
		return errors.New("--user not defined, see usage")
	}

	if config.Password == "" {
		return errors.New("--pass not defined, see usage")
	}

	return nil
}

func performUpdate(c Config) error {
	// Build the url
	updateURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return errors.Wrap(err, "cannot create url")
	}

	// Add the host query parameter to the update url.
	updateParams := url.Values{}
	updateParams.Add("host", c.Host)

	updateURL.RawQuery = updateParams.Encode()

	req, err := http.NewRequest("GET", updateURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	req.SetBasicAuth(c.User, c.Password)

	log.Printf("Updating record %s...", c.Host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot perform http get")
	}
	defer res.Body.Close()

	// Log based on request status code
	switch res.StatusCode {
	case http.StatusOK:
		log.Printf("update of %s was successful.\n", c.Host)
	case http.StatusBadRequest:
		return errors.New("failed to perform request due to invalid input parameters")
	case http.StatusUnauthorized:
		return errors.New("failed to perform request due to authentication issues")
	case http.StatusNotFound:
		return errors.New("failed to perform request because the host you'd like to update cannot be found")
	default:
		return fmt.Errorf("some unknown error occurred: %v", res.Status)
	}

	return nil
}

func main() {
	base := flag.String("base", "https://ydns.io/api/v1/update/", "Base url for api calls on ydns")
	host := flag.String("host", "", "Host to update")
	user := flag.String("user", "", "API Username for authentication on ynds")
	pass := flag.String("pass", "", "API Password for authentication on ynds")
	daemon := flag.Bool("daemon", false, "Enables the updater as a daemon")
	freq := flag.Int("frequency", 60, "Minutes between updates while in daemon mode")

	flag.Parse()

	config := Config{
		BaseURL:  *base,
		Host:     *host,
		User:     *user,
		Password: *pass,
	}

	// Verify config.
	if err := verifyConfig(config); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	if *daemon {
		for {
			// Perform update and then wait for the set duration of minutes.
			if err := performUpdate(config); err != nil {
				log.Fatal(err.Error())
			}

			log.Printf("Now waiting %d minutes.", *freq)
			time.Sleep(time.Duration(*freq) * time.Minute)
		}
	}

	if err := performUpdate(config); err != nil {
		log.Fatal(err.Error())
	}
}
