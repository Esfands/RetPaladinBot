package bot

import (
	"fmt"

	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

func OnUserNoticeMessage(conn *twitch.Client, message twitch.UserNoticeMessage) {
	switch message.MsgID {
	case "sub", "resub":
		var emotes []string = []string{
			"Pog",
			"PogU",
			"POGGERS",
			"PagChomp",
			"PagMan",
			"PagBounce",
		}

		conn.Say(message.Channel, fmt.Sprintf("@EsfandTV, %v %v", message.SystemMsg, utils.GetRandomStringFromSlice(emotes)))
	default:
		fmt.Println("Unknown subscription type")
	}
}
