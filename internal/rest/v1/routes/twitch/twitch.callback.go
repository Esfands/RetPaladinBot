package twitch

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/internal/services/auth"
	"github.com/esfands/retpaladinbot/pkg/errors"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/nicklaw5/helix/v2"
)

func (rg *RouteGroup) LoginCallback(ctx *respond.Ctx) error {
	queries := ctx.Queries()
	state := queries["state"]
	code := queries["code"]

	// Ensure state and code is present
	if state == "" || code == "" {
		slog.Error("[twitch-login-callback] missing state or code", "state", state, "code", code)
		return errors.ErrUnauthorized()
	}

	// Validate the CSRF token
	csrfToken, _, err := rg.gctx.Crate().Auth.ValidateCSRFToken(
		state,
		utils.B2S(ctx.Request().Header.Cookie(auth.CookieCSRF)),
	)
	if err != nil {
		slog.Error("[twitch-login-callback] CSRF validation failed", "state", state, "error", err)
		return errors.ErrUnauthorized()
	}

	ctx.Cookie(csrfToken)

	// Exchange code for token
	twitchToken, err := rg.gctx.Crate().Auth.TwitchExchange(
		context.Background(),
		code,
	)
	if err != nil {
		slog.Error("[twitch-login-callback] error exchanging code", "error", err.Error())
		return err
	}

	// Set the user access token
	rg.gctx.Crate().Helix.Client().SetUserAccessToken(twitchToken.AccessToken)

	// Get user that authenticated
	userReq, err := rg.gctx.Crate().Helix.Client().GetUsers(&helix.UsersParams{})
	if err != nil {
		slog.Error("[twitch-login-callback] error getting user", "error", err.Error())
		return errors.ErrInternalServerError()
	}

	fmt.Println(userReq.Data.Users)

	return nil
}
