package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
)

func listAbuseReports(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	reports, err := services.ListAbuseReport(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(reports)
}

func getAbuseReport(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	report, err := services.GetAbuseReport(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(report)
}

func updateAbuseReportStatus(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "DealAbuseReport", true); err != nil {
		return err
	}

	var data struct {
		Status  string `json:"status" validate:"required"`
		Message string `json:"message" validate:"required,max=4096"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	id, _ := c.ParamsInt("id")

	if err := services.UpdateAbuseReportStatus(uint(id), data.Status, data.Message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func createAbuseReport(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var data struct {
		Resource string `json:"resource" validate:"required"`
		Reason   string `json:"reason" validate:"required,max=4096"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	report, err := services.NewAbuseReport(data.Resource, data.Reason, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(report)
}
