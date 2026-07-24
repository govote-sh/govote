package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/govote-sh/govote/internal/address"
	"github.com/govote-sh/govote/internal/secrets"
	"github.com/govote-sh/govote/internal/utils"
)

const baseURL = "https://www.googleapis.com/civicinfo/v2/voterinfo"

func CheckServer(addr address.InputAddress) tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}

	apiKey, err := secrets.GetAPIKey()
	if err != nil {
		return utils.ErrMsg{Err: err}
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return utils.ErrMsg{Err: fmt.Errorf("could not parse baseURL")}
	}

	// Query params
	params := url.Values{}
	params.Add("address", addr.String())
	base.RawQuery = params.Encode()

	// Perform the HTTP GET request. The API key goes in a header, never the
	// URL: a *url.Error stringifies with the full request URL, so a key in the
	// query string would leak into logs and user-visible error messages.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, base.String(), nil)
	if err != nil {
		return utils.ErrMsg{Err: err}
	}
	req.Header.Set("X-Goog-Api-Key", apiKey)
	res, err := c.Do(req)
	if err != nil {
		// Log the unwrapped error; the *url.Error carries the address-bearing URL.
		loggedErr := err
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			loggedErr = urlErr.Err
		}
		log.Error("Could not perform HTTP GET request", "error", loggedErr)
		// Return a generic message: SSH users see this verbatim.
		return utils.ErrMsg{Err: errors.New("could not reach the election information service")}
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("error closing response body", "error", err)
		}
	}()

	// Check for non-200 response codes
	if res.StatusCode != http.StatusOK {
		return utils.ErrMsg{
			Err:            fmt.Errorf("received non-200 response: %s", res.Status),
			HTTPStatusCode: res.StatusCode,
		}
	}

	// Read and parse the JSON response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrMsg{Err: err, HTTPStatusCode: res.StatusCode}
	}

	// Parse the JSON response into the defined struct
	var data VoterInfoResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return utils.ErrMsg{Err: err, HTTPStatusCode: res.StatusCode}
	}

	// Check if the election day is present
	electionDay := data.Election.ElectionDay
	if electionDay == "" {
		return utils.ErrMsg{Err: fmt.Errorf("could not extract election day from response"), HTTPStatusCode: res.StatusCode}
	}

	return data
}
