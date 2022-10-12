package config

type Config struct {
	AppID string `envconfig:"YADA_APP_ID"`
	Token string `envconfig:"YADA_TOKEN"`

	GuildID string `envconfig:"YADA_GUILD_ID"`

	ImagesChannelID string `envconfig:"YADA_IMAGES_CHANNEL_ID"`
	MusesChannelID  string `envconfig:"YADA_MUSES_CHANNEL_ID"`

	SoundsDataPath string `envconfig:"YADA_SOUNDS_PATH" default:"data/sounds"`

	TenorAPIKey string `envconfig:"TENOR_API_KEY"`
}
