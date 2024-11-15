package dadjoke

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

func NewDadJokeCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "dadjoke"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Get a dad joke 4Head"
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gives you a classic dad joke.",
		"<br/>",
		fmt.Sprintf("<code>%vdadjoke</code>", prefix),
	}
}

func (c *Command) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
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

	url := "https://icanhazdadjoke.com/"

	s := sling.New().Base(url).Set("Accept", "application/json")
	req, err := s.New().Get("/").Request()
	if err != nil {
		slog.Error("Failed to create dad joke request", "error", err)
		return fmt.Sprintf("@%v failed to get a dad joke. FeelsBadMan", user.Name), err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Failed to get dad joke", "error", err)
		return fmt.Sprintf("@%v failed to get a dad joke. FeelsBadMan", user.Name), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read dad joke response body", "error", err)
		return fmt.Sprintf("@%v failed to get a dad joke. FeelsBadMan", user.Name), err
	}

	var joke Response
	if err := json.Unmarshal(body, &joke); err != nil {
		slog.Error("Failed to unmarshal dad joke response", "error", err)
		return fmt.Sprintf("@%v failed to get a dad joke. FeelsBadMan", user.Name), err
	}

	return fmt.Sprintf("@%v %v", target, joke.Joke), nil
}
