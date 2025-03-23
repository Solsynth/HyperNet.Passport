package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
)

func listProgram(c *fiber.Ctx) error {
	var programs []models.Program
	if err := database.C.Find(&programs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(programs)
}

func getProgram(c *fiber.Ctx) error {
	var program models.Program
	programId, _ := c.ParamsInt("programId")
	if err := database.C.Where("id = ?", programId).First(&program).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.JSON(program)
}

func listProgramMembership(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	var members []models.ProgramMember
	if err := database.C.Where("account_id = ?", user.ID).Preload("Program").Find(&members).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(members)
}

func joinProgram(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	programId, _ := c.ParamsInt("programId")
	var program models.Program
	if err := database.C.Where("id = ?", programId).First(&program).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if member, err := services.JoinProgram(user, program); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(member)
	}
}

func leaveProgram(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	programId, _ := c.ParamsInt("programId")
	var program models.Program
	if err := database.C.Where("id = ?", programId).First(&program).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err := services.LeaveProgram(user, program); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
