package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/govote-sh/govote/internal/secrets"
	utils "github.com/govote-sh/govote/internal/utils"
)

const baseURL = "https://www.googleapis.com/civicinfo/v2/voterinfo"

func CheckServer(address string) tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}

	apiKey, err := secrets.GetAPIKey()
	if err != nil {
		return utils.ErrMsg{Err: err}
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return utils.ErrMsg{Err: fmt.Errorf("Could not parse baseURL")}
	}

	// Query params
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("address", address)
	base.RawQuery = params.Encode()

	res, err := c.Get(base.String())

	if err != nil {
		return utils.ErrMsg{Err: err}
	}
	defer res.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrMsg{Err: err}
	}

	// Parse the JSON response into the defined struct
	var data VoterInfoResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return utils.ErrMsg{Err: err}
	}

	// REFACTOR: with better check of API response
	electionDay := data.Election.ElectionDay
	if electionDay == "" {
		return utils.ErrMsg{Err: fmt.Errorf("could not extract election day from response")}
	}

	// Return the election day as a ResData message
	return data
}
