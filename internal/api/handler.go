package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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
	params.Add("key", apiKey)
	params.Add("address", addr.String())
	base.RawQuery = params.Encode()

	// Perform the HTTP GET request
	res, err := c.Get(base.String())
	if err != nil {
		log.Error("Could not perform HTTP GET request", "error", err)
		return utils.ErrMsg{Err: err}
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
