package ydns

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Run will run the update operation.
func Run(base, host, user, pass string) error {
	u, err := url.Parse(base)
	if err != nil {
		return errors.Wrap(err, "cannot create url")
	}

	values := url.Values{}
	values.Set("host", host)

	u.RawQuery = values.Encode()

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read the body")
	}

	logrus.WithFields(logrus.Fields{
		"body":   string(body),
		"status": res.StatusCode,
	}).Debug("got response from api")

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
