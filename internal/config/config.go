package config

type Config struct {
	AppID string `envconfig:"YADA_APP_ID"`
	Token string `envconfig:"YADA_TOKEN"`

	GuildID string `envconfig:"YADA_GUILD_ID"`

	ImagesChannelID string `envconfig:"YADA_IMAGES_CHANNEL_ID"`
	MusesChannelID  string `envconfig:"YADA_MUSES_CHANNEL_ID"`

	VitalyUserID string `envconfig:"YADA_VITALY_USER_ID"`
	LezhikUserID string `envconfig:"YADA_LEZHIK_USER_ID"`
	OlegUserID   string `envconfig:"YADA_OLEG_USER_ID"`
	VeraUserID   string `envconfig:"YADA_VERA_USER_ID"`
}
