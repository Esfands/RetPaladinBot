package song

import (
	"encoding/json"
	"fmt"
	"io"
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

func NewSongCommand(gctx global.Context) *Command {
	cmd := &Command{
		gctx: gctx,
	}

	return cmd
}

func (c *Command) Name() string {
	return "song"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Get the latest track Esfand has listened to."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gets the current song playing in the stream.",
		"<br/>",
		fmt.Sprintf("<code>%vsong</code>", prefix),
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

	req, err := sling.New().Get(fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=esfandtv&api_key=%v&format=json", c.gctx.Config().APIKeys.LastFM)).Request()
	if err != nil {
		return "", err
	}

	request, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	var history Response
	err = json.Unmarshal(body, &history)
	if err != nil {
		return "", err
	}

	if len(history.RecentTracks.Track) == 0 {
		return fmt.Sprintf("@%v, nothing has been listened to yet.", user.Name), nil
	}

	return fmt.Sprintf(
		"@%v, current song: %v - %v | Full history -> https://www.last.fm/user/esfandtv/library",
		target,
		history.RecentTracks.Track[0].Name,
		history.RecentTracks.Track[0].Artist.Text,
	), nil
}
