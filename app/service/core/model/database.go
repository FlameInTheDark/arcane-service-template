package model

import "time"

type GuildCommand struct {
	GuildID string
	ID      string
	Active  bool
}

type GuildActions struct {
	UserID    string
	GuildID   string
	Command   string
	Timestamp time.Time
}
