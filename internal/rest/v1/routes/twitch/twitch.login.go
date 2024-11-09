package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/internal/services/auth"
	"github.com/esfands/retpaladinbot/pkg/errors"
	"github.com/nicklaw5/helix/v2"
)

func (rg *RouteGroup) Login(ctx *respond.Ctx) error {
	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		slog.Error("[twitch-login] error generating bytes", "error", err.Error())
		return errors.ErrInternalServerError()
	}

	state := hex.EncodeToString(tokenBytes[:])

	// Create CSRF token
	csrfToken, err := rg.gctx.Crate().Auth.CreateCSRFToken(state)
	if err != nil {
		slog.Error("[twitch-login] error creating csrf token", "error", err.Error())
		return errors.ErrInternalServerError()
	}

	// Set state cookie
	cookie := rg.gctx.Crate().Auth.Cookie(auth.CookieCSRF, csrfToken, time.Minute*5)
	ctx.Cookie(cookie)

	// Redirect to provider
	return ctx.Redirect(rg.gctx.Crate().Helix.Client().GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       rg.gctx.Crate().Auth.GetTwitchScopes(),
		State:        state,
		ForceVerify:  false,
	}), http.StatusTemporaryRedirect)
}
