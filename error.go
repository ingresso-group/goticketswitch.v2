package ticketswitch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// F13AuthErrorCode is the error code that gets returned for Authentication errors
const F13AuthErrorCode int = 3

// Error represents an error returned by the API
type Error struct {
	Code                int    `json:"error_code"`
	Description         string `json:"error_desc"`
	AuthenticationError bool
	CallbackGoneError   bool
}

func (err Error) Error() string {
	return fmt.Sprintf("ticketswitch: API error %d: %s", err.Code, err.Description)
}

func checkForError(resp *http.Response) error {
	var ret Error
	if resp.StatusCode == http.StatusGone {
		ret.CallbackGoneError = true
	}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&ret)
	if err != nil {
		return err
	}

	if ret.Code == F13AuthErrorCode {
		ret.AuthenticationError = true
	}
	if ret.Code > 0 || ret.Description != "" || ret.AuthenticationError || ret.CallbackGoneError {
		return ret
	}
	return nil
}
