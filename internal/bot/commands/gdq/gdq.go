package gdq

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

type Command struct {
	gctx global.Context
}

func NewGDQCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "gamesdonequick"
}

func (c *Command) Aliases() []string {
	return []string{"gdq"}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Get a random stream donation from any GDQ event."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gives you a random donation from the GDQ events",
		"<br/>",
		fmt.Sprintf("<code>%vgdq</code>", prefix),
		"<br/>",
		fmt.Sprintf("<code>%vgamesdonequick</code>", prefix),
	}
}

func (c *Command) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  false,
		EnabledOffline: true,
	}
}

func (c *Command) UserCooldown() int {
	return 30
}

func (c *Command) GlobalCooldown() int {
	return 10
}

func (c *Command) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	url := "https://api.ivr.fi/"
	s := sling.New().Base(url).Set("Accept", "application/json")
	req, err := s.New().Get("v2/misc/gdq/random").Request()
	if err != nil {
		slog.Error("[gdq-cmd] error getting gdq donation", "error", err.Error())
		return "Error getting GDQ donation FeelsBadMan", nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[gdq-cmd] error getting gdq donation", "error", err.Error())
		return "Error getting GDQ donation FeelsBadMan", nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[gdq-cmd] error reading gdq donation body", "error", err.Error())
		return "Error reading GDQ donation FeelsBadMan", nil
	}

	var gdqResp Response
	err = json.Unmarshal(body, &gdqResp)
	if err != nil {
		slog.Error("[gdq-cmd] error unmarshalling gdq donation", "error", err.Error())
		return "Error unmarshalling GDQ donation FeelsBadMan", nil
	}

	return fmt.Sprintf("@%v [%v] %v", target, gdqResp.EventName, gdqResp.Comment), nil
}
