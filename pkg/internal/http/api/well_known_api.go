package api

import (
	"fmt"

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
		"id_token_signing_alg_values_supported":            []string{"HS512"},
		"token_endpoint_auth_signing_alg_values_supported": []string{"HS512"},
	})
}
