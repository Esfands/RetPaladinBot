package domain

import (
	"github.com/gempir/go-twitch-irc/v4"
)

type CustomCommand struct {
	Name       string `json:"name"`
	Response   string `json:"response"`
	UsageCount int    `json:"usage_count"`
}

type Command struct {
	Name               string       `json:"name"`
	Aliases            []string     `json:"aliases"`
	Permissions        []Permission `json:"permissions"`
	Description        string       `json:"description"`
	DynamicDescription []string     `json:"dynamic_description"`
	GlobalCooldown     int          `json:"global_cooldown"`
	UserCooldown       int          `json:"user_cooldown"`
	EnabledOffline     bool         `json:"enabled_offline"`
	EnabledOnline      bool         `json:"enabled_online"`
	UsageCount         int          `json:"usage_count"`
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
