package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
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

func performUpdate(c Config) {

	// Build the url
	updateURL, err := url.Parse(c.BaseURL)
	if err != nil {
		log.Fatal("Cannot create url:", err)
	}

	// Add the host query parameter to the update url.
	updateParams := url.Values{}
	updateParams.Add("host", c.Host)

	updateURL.RawQuery = updateParams.Encode()

	req, err := http.NewRequest("GET", updateURL.String(), nil)
	if err != nil {
		log.Fatal("Cannot create request:", err)
	}

	req.SetBasicAuth(c.User, c.Password)

	log.Printf("Updating record %s...", c.Host)

	resp, err := http.DefaultClient.Do(req)
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

	base := flag.String("base", "https://ydns.io/api/v1/update/", "Base url for api calls on ydns")
	host := flag.String("host", "", "Host to update")
	user := flag.String("user", "", "API Username for authentication on ynds")
	pass := flag.String("pass", "", "API Password for authentication on ynds")
	daemon := flag.Bool("daemon, d", false, "Enables the updater as a daemon")
	freq := flag.Int("frequency, f", 60, "Minutes inbetween updates while in daemon mode")

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
			performUpdate(config)

			log.Printf("Now waiting %d minutes.", *freq)
			time.Sleep(time.Duration(*freq) * time.Minute)
		}
	} else {

		// Perform update now.
		performUpdate(config)
	}
}
