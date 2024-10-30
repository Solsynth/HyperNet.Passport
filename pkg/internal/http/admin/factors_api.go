package admin

import (
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
)

func getUserAuthFactors(c *fiber.Ctx) error {
	userId, _ := c.ParamsInt("user")

	if err := exts.EnsureGrantedPerm(c, "AdminAuthFactors", true); err != nil {
		return err
	}

	var factors []models.AuthFactor
	if err := database.C.Where("account_id = ?", userId).Find(&factors).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	encodedResp := lo.Map(factors, func(item models.AuthFactor, idx int) map[string]any {
		var encoded map[string]any
		raw, _ := jsoniter.Marshal(item)
		_ = jsoniter.Unmarshal(raw, &encoded)

		// Blur out the secret if it isn't current rolling email one-time-password
		if item.Type != models.EmailPasswordFactor && len(item.Secret) != 6 {
			encoded["secret"] = "**CENSORED**"
		} else {
			encoded["secret"] = item.Secret
		}

		return encoded
	})

	return c.JSON(encodedResp)
}
