package ticketswitch

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	e := Error{
		Code:        123,
		Description: "hampster dead",
	}
	assert.Equal(t, "ticketswitch: API error 123: hampster dead", e.Error())
}

func TestCheckForError(t *testing.T) {
	data := []byte(`{
	"foo": "bar",
	"lol": "beans",
	"code": 123,
	"desc": "some cool stuff"
}`)
	responseWriter := httptest.NewRecorder()
	responseWriter.Write(data)
	response := responseWriter.Result()
	err := checkForError(response)
	assert.Nil(t, err)

	data = []byte(`{
	"error_code": 123,
	"error_desc": "hampster dead"
}`)
	responseWriter = httptest.NewRecorder()
	responseWriter.Write(data)
	response = responseWriter.Result()
	err = checkForError(response)
	assert.NotNil(t, err)
	assert.IsType(t, Error{}, err)
}
