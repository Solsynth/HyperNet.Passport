package web

import (
	"strings"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit"

	"git.solsynth.dev/hypernet/passport/pkg/internal/web/admin"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/api"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type App struct {
	app *fiber.App
}

var IReader *sec.InternalTokenReader

func NewServer() *App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "HyperNet.Passport",
		AppName:               "HyperNet.Passport",
		ReadBufferSize:        5 * 1024 * 1024, // 5MB for large JWT
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
	})

	app.Use(idempotency.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	app.Use(sec.ContextMiddleware(IReader))
	app.Use(authkit.ParseAccountMiddleware)

	admin.MapControllers(app, "/api/admin")
	api.MapControllers(app, "/api")

	return &App{app}
}

func (v *App) Listen() {
	if err := v.app.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting http...")
	}
}
