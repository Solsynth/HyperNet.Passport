package services

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var EReader *sec.JwtReader
var EWriter *sec.JwtWriter

type PayloadClaims struct {
	jwt.RegisteredClaims

	// Internal Stuff
	SessionID string `json:"sed"`

	// ID Token Stuff
	Name  string `json:"name,omitempty"`
	Nick  string `json:"preferred_username,omitempty"`
	Email string `json:"email,omitempty"`

	// Additional Stuff
	AuthorizedParties string `json:"azp,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	Type              string `json:"typ"`
}

const (
	JwtAccessType  = "access"
	JwtRefreshType = "refresh"
)

func EncodeJwt(id string, typ, sub, sed string, nonce *string, aud []string, exp time.Time, idTokenUser ...models.Account) (string, error) {
	var azp string
	for _, item := range aud {
		if item != InternalTokenAudience {
			azp = item
			break
		}
	}

	claims := PayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Audience:  aud,
			Issuer:    viper.GetString("security.issuer"),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
		AuthorizedParties: azp,
		SessionID:         sed,
		Type:              typ,
	}

	if len(idTokenUser) > 0 {
		user := idTokenUser[0]
		claims.Name = user.Name
		claims.Nick = user.Nick
		claims.Email = user.GetPrimaryEmail().Content
	}

	if nonce != nil {
		claims.Nonce = *nonce
	}

	return sec.WriteJwt(EWriter, claims)
}
