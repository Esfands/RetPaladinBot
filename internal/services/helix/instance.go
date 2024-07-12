package helix

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

type Service interface {
	Client() *helix.Client
}

type helixService struct {
	client *helix.Client
}

func (h *helixService) Client() *helix.Client {
	return h.client
}

func (h *helixService) refreshAppAccessToken() error {
	res, err := h.Client().RequestAppAccessToken(
		[]string{"user:read:email"},
	)
	if err != nil {
		return err
	}

	h.Client().SetAppAccessToken(res.Data.AccessToken)

	fmt.Println("twitch app access token refreshed", res.Data.ExpiresIn)
	return nil
}
