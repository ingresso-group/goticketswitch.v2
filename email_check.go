package ticketswitch

// EmailCheckResult represents a result returned
// form check-email endpoint.
type EmailCheckResult struct {
	Code string `json:"error_code,omitempty"`
	Key  string `json:"error_key,omitempty"`
	Desc string `json:"error_desc,omitempty"`
}

func (r *EmailCheckResult) emailIsValid() bool {
	return r.Code == ""
}

// EmailCheckParams represents payload to be sent
// to email-check endpoint to validate given email.
type EmailCheckParams struct {
	UniversalParams
	EmailAddress string
}

func (params *EmailCheckParams) Params() map[string]string {
	values := make(map[string]string)

	for k, v := range params.Universal() {
		values[k] = v
	}

	return values

}
