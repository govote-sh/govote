package api

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"charm.land/log/v2"
	"github.com/govote-sh/govote/internal/address"
	"github.com/govote-sh/govote/internal/secrets"
	"github.com/govote-sh/govote/internal/utils"
)

const testAPIKey = "AIzaSyFAKE-SECRET-KEY-DO-NOT-LEAK"

// TestCheckServerDoesNotLeakAPIKeyOnTransportError verifies that when the HTTP
// request fails at the transport layer, the returned error does not expose the
// API key. The key is sent in a header rather than the URL precisely so that
// the *url.Error the client produces (which embeds the full request URL in its
// message) can never contain it.
func TestCheckServerDoesNotLeakAPIKeyOnTransportError(t *testing.T) {
	t.Setenv("API_KEY", testAPIKey)
	if err := secrets.SetupSecrets(); err != nil {
		t.Fatalf("SetupSecrets: %v", err)
	}
	// Force every outbound request to fail at connect time.
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:9")

	msg := CheckServer(address.InputAddress{Street: "1234 W Broad St", City: "Richmond", State: "VA"})

	errMsg, ok := msg.(utils.ErrMsg)
	if !ok {
		t.Fatalf("expected utils.ErrMsg, got %T: %v", msg, msg)
	}
	if strings.Contains(errMsg.Error(), testAPIKey) {
		t.Fatalf("API key leaked in error message: %q", errMsg.Error())
	}
}

// TestCheckServerDoesNotLogAddressOnTransportError verifies the address (which
// rides in the request URL) is not logged when the transport fails: the raw
// *url.Error must be unwrapped before logging.
func TestCheckServerDoesNotLogAddressOnTransportError(t *testing.T) {
	t.Setenv("API_KEY", testAPIKey)
	if err := secrets.SetupSecrets(); err != nil {
		t.Fatalf("SetupSecrets: %v", err)
	}
	// Force every outbound request to fail at connect time.
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:9")

	// Capture log output; the default logger is a shared global, so restore it.
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	const (
		street = "1234 W Broad St"
		city   = "Richmond"
	)
	msg := CheckServer(address.InputAddress{Street: street, City: city, State: "VA"})

	if _, ok := msg.(utils.ErrMsg); !ok {
		t.Fatalf("expected utils.ErrMsg, got %T: %v", msg, msg)
	}

	logged := buf.String()
	if logged == "" {
		t.Fatal("expected a transport error to be logged, got no log output")
	}
	if strings.Contains(logged, street) || strings.Contains(logged, city) {
		t.Fatalf("user address leaked in server logs: %q", logged)
	}
	if strings.Contains(logged, baseURL) {
		t.Fatalf("request URL leaked in server logs: %q", logged)
	}
}
