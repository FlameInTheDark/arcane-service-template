package model

type Command struct {
	UserID    string `json:"userid"`
	GuildID   string `json:"guildid"`
	ChannelID string `json:"channelid"`
	MessageID string `json:"messageid"`
	Username  string `json:"username"`
	Content   string `json:"content"`
}

type RegisterCommand struct {
	Name   string
	Worker string
}
