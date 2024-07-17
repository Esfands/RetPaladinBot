package variables

import (
	"context"
	"strings"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
)

type VariableI interface {
	GetName() string
	GetAliases() []string
	Code(
		ctx context.Context,
		user twitch.User,
		cmdContext []string,
		context []string,
		message string,
	) string
}

type Variable struct {
	// Name is the name of the variable
	Name string `json:"name"`
}

type Service struct {
	Variables []VariableI
}

type ServiceI interface {
	ParseVariables(
		ctx context.Context,
		user twitch.User,
		cmdContext []string,
		context []string,
	) string
}

// Adds a variable to the list of variables
func (s *Service) registerVariable(variable VariableI) {
	s.Variables = append(s.Variables, variable)
}

func NewService(gctx global.Context) ServiceI {
	svc := &Service{}

	svc.registerVariable(NewUserVariable(gctx))

	return svc
}

func (s *Service) ParseVariables(ctx context.Context, user twitch.User, cmdContext []string, context []string) string {
	// Join the context into a single message string
	message := strings.Join(context, " ")

	// Iterate over the registered variables
	for _, variable := range s.Variables {
		// Replace occurrences of the variable in the message
		placeholder := "${" + variable.GetName() + "}"
		message = strings.ReplaceAll(message, placeholder, variable.Code(ctx, user, cmdContext, context, message))
	}

	return message
}
