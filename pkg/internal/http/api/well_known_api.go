package api

import (
	"fmt"

	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func getOidcConfiguration(c *fiber.Ctx) error {
	domain := viper.GetString("domain")
	basepath := fmt.Sprintf("https://%s", domain)

	return c.JSON(fiber.Map{
		"issuer":                                           viper.GetString("security.issuer"),
		"authorization_endpoint":                           fmt.Sprintf("%s/authorize", basepath),
		"token_endpoint":                                   fmt.Sprintf("%s/api/auth/token", basepath),
		"userinfo_endpoint":                                fmt.Sprintf("%s/api/users/me", basepath),
		"response_types_supported":                         []string{"code", "token"},
		"grant_types_supported":                            []string{"authorization_code", "implicit", "refresh_token"},
		"subject_types_supported":                          []string{"public"},
		"token_endpoint_auth_methods_supported":            []string{"client_secret_post"},
		"id_token_signing_alg_values_supported":            []string{"RS256"},
		"token_endpoint_auth_signing_alg_values_supported": []string{"RS256"},
		"jwks_uri":                                         fmt.Sprintf("%s/.well-known/jwks", basepath),
	})
}

func getJwk(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"keys": []fiber.Map{
			services.EReader.BuildJwk(viper.GetString("id")),
		},
	})
}
