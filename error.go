package ticketswitch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an error returned by the API
type Error struct {
	Code        int    `json:"error_code"`
	Description string `json:"error_desc"`
}

func (err Error) Error() string {
	return fmt.Sprintf("ticketswitch: API error %d: %s", err.Code, err.Description)
}

func checkForError(resp *http.Response) error {
	var ret Error
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&ret)
	if err != nil {
		return err
	}

	if ret.Code > 0 || ret.Description != "" {
		return ret
	}
	return nil
}
