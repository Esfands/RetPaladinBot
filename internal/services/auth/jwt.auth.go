package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (a *authmen) SignJWT(secret string, claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Sign it
	tokenStr, err := token.SignedString(utils.S2B(secret))

	return tokenStr, err
}

func (a *authmen) VerifyJWT(token []string, out jwt.Claims) (*jwt.Token, error) {
	result, err := jwt.ParseWithClaims(
		strings.Join(token, "."),
		out,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("bad jwt signing method, expected HMAC but got %v", t.Header["alg"])
			}

			return utils.S2B(a.JWTSecret), nil
		},
	)

	return result, err
}

type JWTClaimUser struct {
	TwitchID string `json:"tid"`

	jwt.RegisteredClaims
}

type JWTClaimOAuth2CSRF struct {
	State     string    `json:"s"`
	CreatedAt time.Time `json:"at"`
	Bind      string    `json:"bind"`

	jwt.RegisteredClaims
}
