package twitch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/esfands/retpaladinbot/internal/db"
	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/pkg/errors"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/exp/slog"
)

func (rg *RouteGroup) EventSubRecievedNotification(ctx *respond.Ctx) error {
	body := ctx.Request().Body()

	// Verify Twitch sent the event
	mac := hmac.New(sha256.New, utils.S2B(rg.gctx.Config().Twitch.Helix.EventSubSecret))

	mac.Write(
		utils.S2B(
			fmt.Sprintf(
				"%s%s%s",
				utils.B2S(ctx.Request().Header.Peek("Twitch-Eventsub-Message-Id")),
				utils.B2S(ctx.Request().Header.Peek("Twitch-Eventsub-Message-Timestamp")),
				utils.B2S(body),
			),
		),
	)
	hmacsha256 := fmt.Sprintf("sha256=%s", hex.EncodeToString(mac.Sum(nil)))

	if hmacsha256 != utils.B2S(ctx.Request().Header.Peek("Twitch-Eventsub-Message-Signature")) {
		slog.Error("[eventsub] invalid signature on subscription")
		return errors.ErrInvalidSignature().SetDetail("No valid signature on subscription")
	}

	var vals utils.EventSubNotification
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if err != nil {
		slog.Error("[eventsub] couldn't decode the eventsub notification", "error", err.Error())
		return errors.ErrBadRequest().SetDetail("Could not decode body")
	}

	messageType := utils.B2S(ctx.Request().Header.Peek("Twitch-Eventsub-Message-Type"))

	fmt.Println(messageType)
	switch messageType {
	case "notification":
		ctx.Response().SetStatusCode(http.StatusOK)

		switch vals.Subscription.Type {
		case "stream.online":
			var streamOnlinePayload helix.EventSubStreamOnlineEvent
			if err := json.Unmarshal(vals.Event, &streamOnlinePayload); err != nil {
				slog.Error("[eventsub] couldn't unmarshal the stream.online event", "error", err.Error())
				return errors.ErrBadRequest().SetDetail("Could not unmarshal stream.online event")
			}

			rg.streamOnline(streamOnlinePayload)

		case "stream.offline":
			var streamOfflinePayload helix.EventSubStreamOfflineEvent
			if err := json.Unmarshal(vals.Event, &streamOfflinePayload); err != nil {
				slog.Error("[eventsub] couldn't unmarshal the stream.offline event", "error", err.Error())
				return errors.ErrBadRequest().SetDetail("Could not unmarshal stream.offline event")
			}

			rg.streamOffline(streamOfflinePayload)

		case "channel.update":
			var channelUpdatePayload helix.EventSubChannelUpdateEvent
			if err := json.Unmarshal(vals.Event, &channelUpdatePayload); err != nil {
				slog.Error("[eventsub] couldn't unmarshal the channel.update event", "error", err.Error())
				return errors.ErrBadRequest().SetDetail("Could not unmarshal channel.update event")
			}

			rg.channelUpdate(channelUpdatePayload)
		}

	case "webhook_callback_verification":
		fmt.Println("=== CHALLENGE ===")
		ctx.Response().SetStatusCode(http.StatusOK)

		_, err = ctx.Write(utils.S2B(vals.Challenge))
		if err != nil {
			slog.Error("[eventsub] couldn't write the challenge", "error", err.Error())
			return errors.ErrInternalServerError().SetDetail("Could not write challenge")
		}

	case "revocation":
		ctx.Response().SetStatusCode(http.StatusOK)

		fmt.Println("=== REVOCATION ===")
		fmt.Println(vals.Subscription)
	}

	return nil
}

func (rg RouteGroup) streamOnline(event helix.EventSubStreamOnlineEvent) {
	fmt.Println("=== STREAM ONLINE ===")
	fmt.Println(event)

	// Get the channel information from Helix
	channelInfoRes, err := rg.gctx.Crate().Helix.Client().GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{event.BroadcasterUserID},
	})
	if err != nil {
		slog.Error("[eventsub] couldn't get the channel information", "error", err.Error())
		return
	}

	channelInfo := channelInfoRes.Data.Channels[0]

	rg.gctx.Crate().Turso.Queries().InsertStream(rg.gctx, db.StreamStatus{
		StreamID:  event.ID,
		GameID:    sql.NullString{String: channelInfo.GameID, Valid: true},
		GameName:  sql.NullString{String: channelInfo.GameName, Valid: true},
		Live:      true,
		Title:     sql.NullString{String: channelInfo.Title, Valid: true},
		StartedAt: event.StartedAt.Format(time.RFC3339),
		EndedAt:   sql.NullString{String: "", Valid: false},
	})
}

func (rg RouteGroup) streamOffline(event helix.EventSubStreamOfflineEvent) {
	fmt.Println("=== STREAM OFFLINE ===")
	fmt.Println(event)

	// First get the stream from the database to get the ID of the current live stream
	currentLiveStream, err := rg.gctx.Crate().Turso.Queries().GetLiveStream(rg.gctx)
	if err != nil {
		slog.Error("[eventsub] couldn't get the currently live stream", "error", err.Error())
		return
	}

	// Update the stream that went offline
	if err := rg.gctx.Crate().Turso.Queries().StreamWentOffline(rg.gctx, currentLiveStream.StreamID, sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true}); err != nil {
		slog.Error("[eventsub] couldn't update the stream that went offline", "error", err.Error())
		return
	}
}

func (rg RouteGroup) channelUpdate(event helix.EventSubChannelUpdateEvent) {
	fmt.Println("=== CHANNEL UPDATE ===")
	// First get the stream from the database to get the ID of the latest stream
	recentStream, err := rg.gctx.Crate().Turso.Queries().GetMostRecentStreamStatus(rg.gctx)
	if err != nil {
		slog.Error("[eventsub] couldn't get the most recent stream status", "error", err.Error())
		return
	}

	err = rg.gctx.Crate().Turso.Queries().UpdateStreamInfo(rg.gctx, db.StreamStatus{
		ID:        recentStream.ID,
		StreamID:  recentStream.StreamID,
		GameID:    sql.NullString{String: event.CategoryID, Valid: true},
		GameName:  sql.NullString{String: event.CategoryName, Valid: true},
		Live:      recentStream.Live,
		Title:     sql.NullString{String: event.Title, Valid: true},
		StartedAt: recentStream.StartedAt,
		EndedAt:   recentStream.EndedAt,
	})
	if err != nil {
		slog.Error("[eventsub] couldn't update the stream info", "error", err.Error())
		return
	}

	fmt.Println(recentStream)
}
