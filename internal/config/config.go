package config

type Config struct {
	AppID           string `envconfig:"YADA_APP_ID"`
	Token           string `envconfig:"YADA_TOKEN"`
	GuildID         string `envconfig:"YADA_GUILD_ID"`
	ImagesChannelID string `envconfig:"YADA_IMAGES_CHANNEL_ID"`
	DatabaseURL     string `envconfig:"DATABASE_URL"`
	Goodreads       Goodreads
}

type Goodreads struct {
	SessionCookie string `envconfig:"GOODREADS_SESSION_COOKIE"`
}
