package api

import (
	"strings"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func tryAuthorizeThirdClient(c *fiber.Ctx) error {
	id := c.Query("client_id")
	redirect := c.Query("redirect_uri")

	if len(id) <= 0 || len(redirect) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request, missing query parameters")
	}

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, strings.Split(redirect, "?")[0]) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid callback url")
	}

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var ticket models.AuthTicket
	if err := database.C.Where(&models.AuthTicket{
		AccountID: user.ID,
		ClientID:  &client.ID,
	}).Where("last_grant_at IS NULL").First(&ticket).Error; err == nil {
		if ticket.ExpiredAt != nil && ticket.ExpiredAt.Unix() < time.Now().Unix() {
			return c.JSON(fiber.Map{
				"client": client,
				"ticket": nil,
			})
		} else {
			ticket, _ = services.RotateTicket(ticket)
		}

		return c.JSON(fiber.Map{
			"client": client,
			"ticket": ticket,
		})
	}

	return c.JSON(fiber.Map{
		"client": client,
		"ticket": nil,
	})
}

func authorizeThirdClient(c *fiber.Ctx) error {
	id := c.Query("client_id")
	response := c.Query("response_type")
	redirect := c.Query("redirect_uri")
	nonce := c.Query("nonce")
	scope := c.Query("scope")
	if len(scope) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request params")
	}

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch response {
	case "code":
		// OAuth Authorization Mode
		ticket, err := services.NewOauthTicket(
			user,
			client,
			strings.Split(scope, " "),
			[]string{services.InternalTokenAudience, client.Alias},
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
			&nonce,
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user.ID, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.JSON(fiber.Map{
				"ticket":       ticket,
				"redirect_uri": redirect,
			})
		}
	case "token":
		// OAuth Implicit Mode
		ticket, err := services.NewOauthTicket(
			user,
			client,
			strings.Split(scope, " "),
			[]string{services.InternalTokenAudience, client.Alias},
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
			&nonce,
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else if access, refresh, err := services.GetToken(ticket); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user.ID, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.JSON(fiber.Map{
				"access_token":  access,
				"refresh_token": refresh,
				"redirect_uri":  redirect,
				"ticket":        ticket,
			})
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported response type")
	}
}
