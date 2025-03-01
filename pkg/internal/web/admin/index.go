package admin

import (
	"github.com/gofiber/fiber/v2"
)

func MapControllers(app *fiber.App, baseURL string) {
	admin := app.Group(baseURL)
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)

		admin.Post("/notify/all", notifyAllUser)
		admin.Post("/notify/:user", notifyOneUser)

		admin.Get("/users", listUser)
		admin.Get("/users/:user", getUser)
		admin.Get("/users/:user/factors", getUserAuthFactors)
		admin.Put("/users/:user/permissions", editUserPermission)
		admin.Post("/users/:user/confirm", forceConfirmAccount)
	}
}
