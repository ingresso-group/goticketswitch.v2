package ticketswitch

// Config defines the credentials used to access the API
type Config struct {
	BaseURL     string
	User        string
	Password    string
	SubUser     string
	Language    string
	CryptoBlock string
	DebugMode   bool
}

// NewConfig returns a pointer to a newly created Config.
func NewConfig(user, password string) *Config {
	return &Config{
		BaseURL:  "https://api.ticketswitch.com",
		User:     user,
		Password: password,
	}
}
