package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"gorm.io/gorm"

	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func getUserInBatch(c *fiber.Ctx) error {
	id := c.Query("id")
	list := strings.Split(id, ",")
	var nameList []string
	numericList := lo.Filter(lo.Map(list, func(str string, i int) int {
		value, err := strconv.Atoi(str)
		if err != nil {
			nameList = append(nameList, str)
			return 0
		}
		return value
	}), func(vak int, idx int) bool {
		return vak > 0
	})

	tx := database.C
	if len(numericList) > 0 {
		tx = tx.Where("id IN ?", numericList)
	}
	if len(nameList) > 0 {
		tx = tx.Or("name IN ?", nameList)
	}
	if len(nameList) == 0 && len(numericList) == 0 {
		return c.JSON([]models.Account{})
	}

	var accounts []models.Account
	if err := tx.
		Preload("Profile").
		Preload("Badges", func(db *gorm.DB) *gorm.DB {
			return db.Order("badges.is_active DESC, badges.type DESC")
		}).
		Find(&accounts).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(accounts)
}

func lookupAccount(c *fiber.Ctx) error {
	probe := c.Query("probe")
	if len(probe) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "lookup probe is required")
	}

	user, err := services.LookupAccount(probe)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(user)
}

func searchAccount(c *fiber.Ctx) error {
	probe := c.Query("probe")
	if len(probe) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "search probe is required")
	}

	users, err := services.SearchAccount(probe)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}

func getUserinfo(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		Preload("Contacts").
		Preload("Badges").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		data.PermNodes = c.Locals("nex_user").(*sec.UserInfo).PermNodes
	}

	var resp fiber.Map
	raw, _ := jsoniter.Marshal(data)
	_ = jsoniter.Unmarshal(raw, &resp)

	return c.JSON(resp)
}

func editUserinfo(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Nick        string            `json:"nick" validate:"required"`
		Description string            `json:"description"`
		FirstName   string            `json:"first_name"`
		LastName    string            `json:"last_name"`
		Location    string            `json:"location"`
		TimeZone    string            `json:"time_zone"`
		Gender      string            `json:"gender"`
		Pronouns    string            `json:"pronouns"`
		Links       map[string]string `json:"links"`
		Birthday    time.Time         `json:"birthday"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	} else {
		data.Nick = strings.TrimSpace(data.Nick)
	}
	if !services.ValidateAccountName(data.Nick, 1, 24) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid account nick, length requires 4 to 24")
	}

	var account models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	links := make(map[string]any)
	for k, v := range data.Links {
		links[k] = v
	}

	account.Nick = data.Nick
	account.Profile.Gender = data.Gender
	account.Profile.Pronouns = data.Pronouns
	account.Profile.Location = data.Location
	account.Profile.TimeZone = data.TimeZone
	account.Profile.Links = links
	account.Profile.Description = data.Description
	account.Profile.FirstName = data.FirstName
	account.Profile.LastName = data.LastName
	account.Profile.Birthday = &data.Birthday

	if err := database.C.Save(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if err := database.C.Save(&account.Profile).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	services.AddEvent(user.ID, "profile.edit", nil, c.IP(), c.Get(fiber.HeaderUserAgent))
	services.InvalidUserAuthCache(account.ID)

	return c.SendStatus(fiber.StatusOK)
}

func updateAccountLanguage(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Language string `json:"language" validate:"required,bcp47_language_tag"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := database.C.Model(&models.Account{}).Where("id = ?", user.ID).
		Updates(&models.Account{Language: data.Language}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	services.AddEvent(user.ID, "profile.edit.language", nil, c.IP(), c.Get(fiber.HeaderUserAgent))
	services.InvalidUserAuthCache(user.ID)

	user.Language = data.Language

	return c.JSON(user)
}

func doRegister(c *fiber.Ctx) error {
	var data struct {
		Name         string `json:"name" validate:"required,lowercase,alphanum,min=4,max=16"`
		Nick         string `json:"nick" validate:"required"`
		Email        string `json:"email" validate:"required,email"`
		Password     string `json:"password" validate:"required,min=4,max=32"`
		Language     string `json:"language" validate:"required,bcp47_language_tag"`
		CaptchaToken string `json:"captcha_token" validate:"required"`
		MagicToken   string `json:"magic_token"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	} else {
		data.Name = strings.TrimSpace(data.Name)
		data.Nick = strings.TrimSpace(data.Nick)
		data.Email = strings.TrimSpace(data.Email)
	}
	if _, err := strconv.Atoi(data.Name); err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid account name, cannot be pure number")
	}
	if !services.ValidateAccountName(data.Nick, 1, 24) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid account nick, length requires 4 to 24")
	}
	if viper.GetBool("use_registration_magic_token") && len(data.MagicToken) <= 0 {
		return fmt.Errorf("missing magic token in request")
	} else if viper.GetBool("use_registration_magic_token") {
		if tk, err := services.ValidateMagicToken(data.MagicToken, models.RegistrationMagicToken); err != nil {
			return err
		} else {
			database.C.Delete(&tk)
		}
	}

	if !gap.Nx.ValidateCaptcha(data.CaptchaToken, c.IP()) {
		return fiber.NewError(fiber.StatusBadRequest, "captcha check failed")
	}

	if user, err := services.CreateAccount(
		data.Name,
		data.Nick,
		data.Email,
		data.Password,
		data.Language,
	); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(user)
	}
}

func doRegisterConfirm(c *fiber.Ctx) error {
	var data struct {
		Code string `json:"code" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := services.ConfirmAccount(data.Code); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func reNotifyRegisterConfirm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var magicToken models.MagicToken
	if err := database.C.Where("account_id = ? AND type = ?", user.ID, models.ConfirmMagicToken).First(&magicToken).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := services.NotifyMagicToken(magicToken); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func requestDeleteAccount(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if err := services.CheckAbleToDeleteAccount(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if err = services.RequestDeleteAccount(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func confirmDeleteAccount(c *fiber.Ctx) error {
	var data struct {
		Code string `json:"code" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := services.ConfirmDeleteAccount(data.Code); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
