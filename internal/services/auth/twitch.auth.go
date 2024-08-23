package auth

import (
	"context"
	"log/slog"

	oidc "github.com/coreos/go-oidc"
	"github.com/esfands/retpaladinbot/config"
	"golang.org/x/oauth2"
)

type TwitchAuth struct {
	StateCallbackKey string
	OauthSessionName string
	OauthTokenKey    string
	Scopes           []string
	Claims           oauth2.AuthCodeOption
	Config           *oauth2.Config
	OidcVerifier     *oidc.IDTokenVerifier
}

// initTwitchProvider initializes the Twitch OAuth2 provider.
func (a *authmen) initTwtchProvider(cfg *config.Config) {
	provider, err := oidc.NewProvider(context.Background(), "https://id.twitch.tv/oauth2")
	if err != nil {
		slog.Error("oidc.NewProvider", "error", err)
	}

	scopes := []string{
		"user:read:email",
		"openid",
	}

	TwitchOauth2Config := &oauth2.Config{
		ClientID:     cfg.Twitch.Helix.ClientID,
		ClientSecret: cfg.Twitch.Helix.ClientSecret,
		Scopes:       scopes,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.Twitch.Helix.RedirectURI,
	}

	TwitchOidcVerifier := provider.Verifier(&oidc.Config{ClientID: cfg.Twitch.Helix.ClientID})

	a.Twitch = TwitchAuth{
		StateCallbackKey: "oauth-state-callback",
		OauthSessionName: CookieCSRF,
		OauthTokenKey:    "oauth-token",
		Scopes:           scopes,
		Claims:           oauth2.SetAuthURLParam("claims", `{"id_token":{"email":null}}`),
		Config:           TwitchOauth2Config,
		OidcVerifier:     TwitchOidcVerifier,
	}
}

func (a *authmen) GetTwitchScopes() []string {
	return a.Twitch.Scopes
}

// TwitchExchange exchanges an authorization code into a token
func (a *authmen) TwitchExchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.Twitch.Config.Exchange(ctx, code)
}
