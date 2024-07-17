package v1

import (
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/internal/rest/v1/routes"
	"github.com/esfands/retpaladinbot/internal/rest/v1/routes/commands"
	"github.com/gofiber/fiber/v2"
)

func ctx(fn func(*respond.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		newCtx := &respond.Ctx{Ctx: c}
		return fn(newCtx)
	}
}

func New(gctx global.Context, router fiber.Router) {
	indexRoute := routes.NewRouteGroup(gctx)
	router.Get("/", indexRoute.Index)

	commandRotues := commands.NewRouteGroup(gctx)
	router.Get("/commands", ctx(commandRotues.GetCommands))
}
