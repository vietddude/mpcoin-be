package config

type GoogleOAuthClient struct {
	ClientID     string `env:"GOOGLE_CLIENT_ID" required:"true"`
	ClientSecret string `env:"GOOGLE_CLIENT_SECRET" required:"true"`
	RedirectURI  string `env:"GOOGLE_REDIRECT_URI" required:"true"`
}
