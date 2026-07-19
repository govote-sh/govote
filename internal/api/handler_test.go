package api

import (
	"strings"
	"testing"

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
