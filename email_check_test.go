package ticketswitch

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
    "errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHTTPClient struct {
	http.Client
}

func TestVerifyingEmail(t *testing.T) {
	table := []struct {
		name               string
		email              string
		returnedStatusCode int
		returnedResponse   string
		expectedResult     error
	}{
		{
			name:               "Successful verification",
			email:              "test@gmail.com",
			returnedStatusCode: 200,
			returnedResponse:   `{}`,
			expectedResult:     nil,
		},
		{
			name:               "Invalid email",
			email:              "test@@@@gmail.com",
			returnedStatusCode: 460,
			returnedResponse:   `{"error_code": 9000, "error_desc": "Failzor"}`,
			expectedResult: Error{
				Code:                9000,
				Description:         "Failzor",
				AuthenticationError: false,
				CallbackGoneError:   false,
			},
		},
		{
			name:               "Email not provided",
			email:              "",
			returnedStatusCode: 460,
			returnedResponse:   `{}`,
			expectedResult: errors.New("No email was provided for verification"),
		},
		{
			name:               "Unhandled exception",
			email:              "test@gmail.com",
			returnedStatusCode: 460,
			returnedResponse:   `{"error_code": 2020, "error_desc": "Core broken!"}`,
			expectedResult: Error{
				Code:                2020,
				Description:         "Core broken!",
				AuthenticationError: false,
				CallbackGoneError:   false,
			},
		},
	}

	for _, test := range table {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if test.returnedStatusCode != 200 {
				http.Error(w, test.returnedResponse, test.returnedStatusCode)
			} else {
				fmt.Fprintln(w, test.returnedResponse)
			}
		}))
		cfg := &Config{
			BaseURL:  ts.URL,
			User:     "test",
			Password: "test",
		}
		client := NewClient(cfg)
		params := &EmailCheckParams{
			EmailAddress: test.email,
		}
		ctx := context.Background()
		t.Run(test.name, func(t *testing.T) {
			email_err := client.EmailCheck(ctx, params)
			assert.Equal(t, test.expectedResult, email_err)
		})
	}
}
