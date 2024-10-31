package services

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"strconv"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func GetToken(ticket models.AuthTicket) (atk, rtk string, err error) {
	if err = ticket.IsAvailable(); err != nil {
		return
	}
	if ticket.AccessToken == nil || ticket.RefreshToken == nil {
		err = fmt.Errorf("unable to encode token, access or refresh token id missing")
		return
	}

	atkDeadline := time.Duration(viper.GetInt64("security.access_token_duration")) * time.Second
	rtkDeadline := time.Duration(viper.GetInt64("security.refresh_token_duration")) * time.Second

	sub := strconv.Itoa(int(ticket.AccountID))
	sed := strconv.Itoa(int(ticket.ID))
	atk, err = EncodeJwt(*ticket.AccessToken, JwtAccessType, sub, sed, nil, ticket.Audiences, time.Now().Add(atkDeadline))
	if err != nil {
		return
	}
	rtk, err = EncodeJwt(*ticket.RefreshToken, JwtRefreshType, sub, sed, nil, ticket.Audiences, time.Now().Add(rtkDeadline))
	if err != nil {
		return
	}

	ticket.LastGrantAt = lo.ToPtr(time.Now())
	database.C.Save(&ticket)

	return
}

func ExchangeToken(token string) (atk, rtk string, err error) {
	var ticket models.AuthTicket
	if err = database.C.Where(models.AuthTicket{GrantToken: &token}).First(&ticket).Error; err != nil {
		return
	} else if ticket.LastGrantAt != nil {
		err = fmt.Errorf("ticket was granted the first token, use refresh token instead")
		return
	} else if len(ticket.Audiences) > 1 {
		err = fmt.Errorf("should use authorization code grant type")
		return
	}

	return GetToken(ticket)
}

func ExchangeOauthToken(clientId, clientSecret, redirectUri, token string) (idk, atk, rtk string, err error) {
	var client models.ThirdClient
	if err = database.C.Where(models.ThirdClient{Alias: clientId}).First(&client).Error; err != nil {
		return
	} else if client.Secret != clientSecret {
		err = fmt.Errorf("invalid client secret")
		return
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, redirectUri) {
		err = fmt.Errorf("invalid redirect uri")
		return
	}

	var ticket models.AuthTicket
	if err = database.C.Where(models.AuthTicket{GrantToken: &token}).First(&ticket).Error; err != nil {
		return
	} else if ticket.LastGrantAt != nil {
		err = fmt.Errorf("ticket was granted the first token, use refresh token instead")
		return
	}

	atk, rtk, err = GetToken(ticket)
	if err != nil {
		return
	}

	var user models.Account
	if err = database.C.Where(models.Account{
		BaseModel: models.BaseModel{ID: ticket.AccountID},
	}).Preload("Contacts").First(&user).Error; err != nil {
		return
	}

	sub := strconv.Itoa(int(ticket.AccountID))
	sed := strconv.Itoa(int(ticket.ID))
	idk, err = EncodeJwt(*ticket.AccessToken, JwtAccessType, sub, sed, ticket.Nonce, ticket.Audiences, time.Now().Add(24*time.Minute), user)

	return
}

func RefreshToken(token string) (atk, rtk string, err error) {
	parseInt := func(str string) int {
		val, _ := strconv.Atoi(str)
		return val
	}

	var ticket models.AuthTicket
	var claims *PayloadClaims
	if claims, err = sec.ReadJwt(EReader, token, &PayloadClaims{}); err != nil {
		return
	}

	if claims.Type != JwtRefreshType {
		err = fmt.Errorf("invalid token type, expected refresh token")
		return
	} else if err = database.C.Where(models.AuthTicket{
		BaseModel: models.BaseModel{ID: uint(parseInt(claims.SessionID))},
	}).First(&ticket).Error; err != nil {
		return
	}

	if ticket, err = RotateTicket(ticket); err != nil {
		return
	} else {
		return GetToken(ticket)
	}
}
