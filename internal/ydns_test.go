package ydns_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ydns "github.com/wyattjoh/ydns-updater/internal"
)

func TestRun(t *testing.T) {
	// This can be run "in parallel".
	t.Parallel()

	var host string

	// Setup the testing server to get the request.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the passed "host" request parameter as the external "host" var.
		host = r.URL.Query().Get("host")

		r.URL.User.Username()

		// Send the OK back.
		fmt.Fprintln(w, "ok")
	}))
	defer ts.Close()

	// Ensure there's no error running it with the right options.
	if err := ydns.Run(ts.URL, "test-1.com", "user", "pass"); err != nil {
		t.Fatalf("expected no errors, got: %v", err)
	}

	if host != "test-1.com" {
		t.Fatalf("expected host to equal test-1.com, got: %s", host)
	}
}
