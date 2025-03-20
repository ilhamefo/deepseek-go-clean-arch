package handler

import (
	"event-registration/internal/constant"
	"event-registration/internal/core/service"
	validate "event-registration/internal/infrastructure/validator"
	"event-registration/internal/request"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ExporterHandler struct {
	service   *service.ExporterService
	validator *validate.Validator
	logger    *zap.Logger
}

func NewExporterHandler(service *service.ExporterService, validator *validate.Validator, logger *zap.Logger) *ExporterHandler {
	return &ExporterHandler{service: service, validator: validator, logger: logger}
}

func (h *ExporterHandler) ExportRekapTransaksi(c *fiber.Ctx) error {
	request := new(request.RekapRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constant.INVALID_REQUEST_BODY,
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error_validations": h.validator.ValidationErrors(err),
		})
	}

	err := h.service.ExportRekapTransaksi(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": constant.SUCCESS_EXPORT})
}

func (h *ExporterHandler) ExportAllRekapTransaksi(c *fiber.Ctx) error {
	request := new(request.RekapRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constant.INVALID_REQUEST_BODY,
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error_validations": h.validator.ValidationErrors(err),
		})
	}

	err := h.service.ExportAllRekapTransaksi(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": constant.SUCCESS_EXPORT})
}

func (h *ExporterHandler) ExportRekapPelanggan(c *fiber.Ctx) error {
	request := new(request.RekapRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constant.INVALID_REQUEST_BODY,
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error_validations": h.validator.ValidationErrors(err),
		})
	}

	err := h.service.ExportRekapPelanggan(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": constant.SUCCESS_EXPORT})
}

func (h *ExporterHandler) HelloWorld(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"hello": "world"})
}
