package domain

import (
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type CustomCommand struct {
	Name     string `json:"name" bson:"name"`
	Response string `json:"response" bson:"response"`
}

type Command struct {
	ID             string     `json:"id"`
	Name           string     `json:"name" bson:"name"`
	Aliases        []string   `json:"aliases" bson:"aliases"`
	Default        bool       `json:"default" bson:"default"`
	Response       string     `json:"response" bson:"response"`
	Enabled        bool       `json:"enabled" bson:"enabled"`
	EnabledOnline  bool       `json:"enabled_online" bson:"enabled_online"`
	EnabledOffline bool       `json:"enabled_offline" bson:"enabled_offline"`
	GlobalCooldown int        `json:"global_cooldown" bson:"global_cooldown"`
	UserCooldown   int        `json:"user_cooldown" bson:"user_cooldown"`
	UpdatedAt      *time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
}

type DefaultCommand interface {
	Name() string
	Aliases() []string
	Permissions() []Permission
	Description() string
	DynamicDescription() []string
	Conditions() DefaultCommandConditions
	GlobalCooldown() int
	UserCooldown() int
	Code(user twitch.User, context []string) (string, error)
}

type DefaultCommandConditions struct {
	EnabledOnline  bool `json:"enabled_online"`
	EnabledOffline bool `json:"enabled_offline"`
}
