package ticketswitch

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancellationResult_IsFullyCancelled(t *testing.T) {
	table := []struct {
		name           string
		testdata       string
		expectedResult bool
	}{
		{
			name:           "Successful cancellation",
			testdata:       "testdata/cancel.json",
			expectedResult: true,
		},
		{
			name:           "Must also cancel response",
			testdata:       "testdata/must_also_cancel.json",
			expectedResult: false,
		},
		{
			name:           "Partial cancellation",
			testdata:       "testdata/partial_cancel.json",
			expectedResult: false,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			testdata, err := ioutil.ReadFile(test.testdata)
			if !assert.Nil(t, err) {
				t.Fatal(err)
			}
			cancellation := &CancellationResult{}
			err = json.Unmarshal(testdata, cancellation)
			if !assert.Nil(t, err) {
				t.Fatal(err)
			}
			assert.Equal(t, test.expectedResult, cancellation.IsFullyCancelled())
		})
	}
}
