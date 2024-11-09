package temperature

import (
	"fmt"
	"strconv"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx global.Context
}

func NewTemperatureCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "temperature"
}

func (c *Command) Aliases() []string {
	return []string{"temp"}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Translates a temperature from Fahrenheit to Celsius and vice versa."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Translate Celsius to Fahrenheit",
		fmt.Sprintf("<code>%vtemperature toF (celsius)</code>", prefix),
		"<br/>",
		"Translate Fahrenheit to Celsius.",
		fmt.Sprintf("<code>%vtemperature toC (fahrenheit)</code>", prefix),
	}
}

func (c *Command) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  false,
		EnabledOffline: true,
	}
}

func (c *Command) UserCooldown() int {
	return 30
}

func (c *Command) GlobalCooldown() int {
	return 10
}

func (c *Command) Code(user twitch.User, context []string) (string, error) {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	if len(context) < 2 {
		return fmt.Sprintf(`Invalid number of arguments. Example: "%vtemperature toF/toC (number)"`, prefix), nil
	}

	switch context[0] {
	case "toF":
		// Convert Celsius to Fahrenheit
		celsius, err := strconv.ParseFloat(context[1], 64)
		if err != nil {
			return "Invalid temperature value. Please provide a valid number.", nil
		}

		fahrenheit := (celsius * 9 / 5) + 32
		return fmt.Sprintf("%v, %.2f°C is %.2f°F", user.Name, celsius, fahrenheit), nil

	case "toC":
		// Convert Fahrenheit to Celsius
		fahrenheit, err := strconv.ParseFloat(context[1], 64)
		if err != nil {
			return "Invalid temperature value. Please provide a valid number.", nil
		}

		celsius := (fahrenheit - 32) * 5 / 9
		return fmt.Sprintf("%v, %.2f°F is %.2f°C", user.Name, fahrenheit, celsius), nil

	default:
		return fmt.Sprintf(`Invalid conversion type. Example: "%vtemperature toF/toC (number)"`, prefix), nil
	}
}
