package ticketswitch

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
    values["email_address"] = params.EmailAddress

	return values

}
