package helix

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/esfands/retpaladinbot/internal/services/scheduler"
	"github.com/nicklaw5/helix/v2"
)

type SetupOptions struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func Setup(ctx context.Context, scheduler scheduler.Service, opts SetupOptions) (Service, error) {
	svc := &helixService{}
	var err error

	svc.client, err = helix.NewClientWithContext(ctx, &helix.Options{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURI:  opts.RedirectURI,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Refresh this token every 50 days
	_, err = scheduler.Scheduler().Every(50).Days().Do(func() {
		slog.Debug("refreshing twitch app access token")
		err := svc.refreshAppAccessToken()
		if err != nil {
			fmt.Println("twitch app access token refresh error", err)
		}
	})
	if err != nil {
		return nil, err
	}

	return svc, nil
}
