package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/esfands/retpaladinbot/config"
	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type Authmen interface {
	SignJWT(secret string, claim jwt.Claims) (string, error)
	VerifyJWT(token []string, out jwt.Claims) (*jwt.Token, error)

	CreateCSRFToken(state string) (token string, err error)
	ValidateCSRFToken(state, cookieData string) (*fiber.Cookie, *JWTClaimOAuth2CSRF, error)
	CreateAccessToken(accountID, integrationID string) (string, time.Time, error)

	Cookie(key, token string, duration time.Duration) *fiber.Cookie

	// Twitch
	GetTwitchScopes() []string
	TwitchExchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type authmen struct {
	// Secret key used for signing.
	JWTSecret string
	// Domain for the cookie.
	Domain string
	// If cookie should be secure or not
	Secure bool

	Twitch TwitchAuth
}

const (
	CookieAuth = "rpb-token"
	CookieCSRF = "rpb-csrf"
)

func Setup(jwtSecret, domain string, secure bool, cfg *config.Config) Authmen {
	a := &authmen{
		JWTSecret: jwtSecret,
		Domain:    domain,
		Secure:    secure,
	}

	a.initTwtchProvider(cfg)

	return a
}

func (a *authmen) CreateCSRFToken(state string) (token string, err error) {
	token, err = a.SignJWT(a.JWTSecret, &JWTClaimOAuth2CSRF{
		State:     state,
		CreatedAt: time.Now(),
	})
	if err != nil {
		slog.Error("csrf_token, sign", "error", err)
		return "", err
	}

	return token, nil
}

func (a *authmen) ValidateCSRFToken(
	state, cookieData string,
) (*fiber.Cookie, *JWTClaimOAuth2CSRF, error) {
	csrfToken := strings.Split(cookieData, ".")

	if len(csrfToken) != 3 {
		return nil, nil, fmt.Errorf(
			"bad state (found %d segments when 3 were expected)",
			len(csrfToken),
		)
	}

	// Verifiy the CSRF token.
	csrfClaim := &JWTClaimOAuth2CSRF{}

	token, err := a.VerifyJWT(csrfToken, csrfClaim)
	if err != nil {
		return nil, csrfClaim, fmt.Errorf("invalid state: %s", err.Error())
	}

	{
		b, err := json.Marshal(token.Claims)
		if err != nil {
			return nil, csrfClaim, fmt.Errorf("invalid state: %s", err.Error())
		}

		if err = json.Unmarshal(b, csrfClaim); err != nil {
			return nil, csrfClaim, fmt.Errorf("invalid state: %s", err.Error())
		}
	}

	// Validate: token date
	if csrfClaim.CreatedAt.Before(time.Now().Add(-time.Minute * 5)) {
		return nil, csrfClaim, fmt.Errorf("expired state")
	}

	// Check mismatch
	if state != csrfClaim.State {
		return nil, csrfClaim, fmt.Errorf("mismatched state value")
	}

	// Udate the CSRF cookie (immediate expire)
	cookie := a.Cookie(CookieCSRF, "", 0)

	return cookie, csrfClaim, nil
}

// CreateAccessToken creates a new access token which represents a user.
func (a *authmen) CreateAccessToken(accountID, integrationID string) (string, time.Time, error) {
	expireAt := time.Now().Add(time.Hour * 24 * 7) // 7 days

	token, err := a.SignJWT(a.JWTSecret, &JWTClaimUser{
		AccountID:     accountID,
		IntegrationID: integrationID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "retpaladinbot-api",
			ExpiresAt: &jwt.NumericDate{Time: expireAt},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	})
	if err != nil {
		slog.Error("access_token, sign", "error", err)
		return "", time.Time{}, err
	}

	return token, expireAt, nil
}

func (a *authmen) Cookie(key, token string, duration time.Duration) *fiber.Cookie {
	cookie := &fiber.Cookie{}
	cookie.Name = key
	cookie.Value = token
	cookie.Expires = time.Now().Add(duration)
	cookie.HTTPOnly = true
	cookie.Domain = a.Domain
	cookie.Path = "/"
	cookie.SameSite = fiber.CookieSameSiteNoneMode

	return cookie
}
