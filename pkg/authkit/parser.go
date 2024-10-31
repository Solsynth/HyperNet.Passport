package authkit

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

// GetAccountFromUserInfo returns the account from the user info
// This method will not to query the database, it will parse the token and get the subject of the userinfo token
func GetAccountFromUserInfo(info *sec.UserInfo) models.Account {
	raw, _ := json.Marshal(info.Metadata)

	// We assume the token is signed by the same version of service
	// So directly read the data out of the metadata
	var out models.Account
	_ = json.Unmarshal(raw, &out)
	return out
}

func ParseAccountMiddleware(c *fiber.Ctx) error {
	if info, ok := c.Locals("nex_user").(*sec.UserInfo); ok {
		c.Locals("user", GetAccountFromUserInfo(info))
	}
	return c.Next()
}
