package variables

import (
	"context"
	"fmt"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type UserVariable struct {
	gctx global.Context

	name string
}

func NewUserVariable(gctx global.Context) VariableI {
	return &UserVariable{
		gctx: gctx,
		name: "user",
	}
}

func (v *UserVariable) GetName() string {
	return v.name
}

func (v *UserVariable) GetAliases() []string {
	return []string{}
}

func (v *UserVariable) Code(_ context.Context, user twitch.User, cmdContext []string, _ []string, _ string) string {
	target := utils.GetTarget(user, cmdContext[1:])

	return fmt.Sprintf("@%v", target)
}
