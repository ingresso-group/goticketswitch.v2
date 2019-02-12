package ticketswitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	// IndentOutput controls wether or not JSON output is indented.
	IndentOutput = true
	// EscapeHTML controls wether json output is HTML escaped.
	EscapeHTML = false
)

// marshal is a custom json marshaller that conditionally turns off html
// escaping and applies indentation.
func marshal(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	e := json.NewEncoder(buf)
	if IndentOutput {
		e.SetIndent("", "  ")
	}
	if EscapeHTML {
		e.SetEscapeHTML(false)
	}

	err := e.Encode(v)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DateRange returns a string representation of a date range between two
// time.Times
func DateRange(from time.Time, to time.Time) string {
	if from.IsZero() && to.IsZero() {
		return ""
	}

	buf := bytes.NewBufferString("")
	if !from.IsZero() {
		buf.WriteString(from.Format("20060102"))
	}
	buf.WriteString(":")
	if !to.IsZero() {
		buf.WriteString(to.Format("20060102"))
	}
	return buf.String()
}

// intArrayToString takes an array on ints and returns them as a comma
// separated string
func intArrayToString(input []int) string {
	stringArray := make([]string, len(input))
	for idx, val := range input {
		stringArray[idx] = fmt.Sprint(val)
	}
	return strings.Join(stringArray, ",")
}
