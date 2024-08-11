package isbanned

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type IsBannedCommand struct {
	gctx global.Context
}

func NewIsBannedCommand(gctx global.Context) *IsBannedCommand {
	return &IsBannedCommand{
		gctx: gctx,
	}
}

func (c *IsBannedCommand) Name() string {
	return "isbanned"
}

func (c *IsBannedCommand) Aliases() []string {
	return []string{}
}

func (c *IsBannedCommand) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *IsBannedCommand) Description() string {
	return "Check if a user is banned on Twitch."
}

func (c *IsBannedCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Check if a user is banned on Twitch.",
		"<br/>",
		fmt.Sprintf("<code>%v!isbanned (username)</code>", prefix),
	}
}

func (c *IsBannedCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *IsBannedCommand) UserCooldown() int {
	return 30
}

func (c *IsBannedCommand) GlobalCooldown() int {
	return 10
}

func (c *IsBannedCommand) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	url := "https://api.ivr.fi"
	s := sling.New().Base(url).Set("Accept", "application/json")
	req, err := s.New().Get(fmt.Sprintf("v2/twitch/user?login=%v", target)).Request()
	if err != nil {
		slog.Error("[isbanned-cmd] failed to create sling request", "error", err)
		return fmt.Sprintf("@%v failed to check if a user is banned. FeelsBadMan", user.Name), err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[isbanned-cmd] Failed to create response with http client", "error", err)
		return fmt.Sprintf("@%v failed to check if a user is banned. FeelsBadMan", user.Name), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[isbanned-cmd] Failed to read is banned response body", "error", err)
		return fmt.Sprintf("@%v failed to check if a user is banned. FeelsBadMan", user.Name), err
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		slog.Error("[isbanned-cmd] Failed to unmarshal ban check response", "error", err)
		return fmt.Sprintf("@%v failed to check if a user is banned. FeelsBadMan", user.Name), err
	}

	if resp.StatusCode != http.StatusOK || len(users) == 0 {
		slog.Error("[isbanned-cmd] The api.ivr.fi failed to check a specific user", "status", resp.StatusCode)
		return fmt.Sprintf("@%v failed to check if a user is banned. FeelsBadMan", user.Name), nil
	} else if len(users) == 0 {
		return fmt.Sprintf("@%v that user doesn't exist... Susage", user.Name), nil
	}

	banCheckUser := users[0]

	if banCheckUser.Banned {
		switch banCheckUser.BanReason {
		case "TOS_INDEFINITE":
			return fmt.Sprintf("@%v, %v is indefinitly banned on Twitch. FeelsBadMan", user.Name, target), nil
		case "DMCA":
			return fmt.Sprintf("@%v, %v is banned on Twitch for violating DMCA. FeelsBadMan GuitarTime", user.Name, target), nil
		case "TOS_TEMPORARY":
			return fmt.Sprintf("@%v, %v is temporarily banned on Twitch. FeelsBadMan", user.Name, target), nil
		default:
			return fmt.Sprintf("@%v, unexpected ban reason: %v", user.Name, banCheckUser.BanReason), nil
		}
	} else {
		return fmt.Sprintf("@%v, %v is not banned on Twitch! PogChamp", user.Name, target), nil
	}
}
