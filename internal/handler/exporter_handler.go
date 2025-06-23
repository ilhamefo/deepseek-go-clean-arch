package handler

import (
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/service"
	validate "event-registration/internal/infrastructure/validator"

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

// Get transaksi by unit id godoc
// @Summary Get transaksi by unit id
// @Description Get transaksi by unit id
// @Tags exporter
// @Accept  json
// @Produce  json
// @Param request body request.RekapRequest false "..."
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string][]string
// @Router /transaksi [post]
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

// Get transaksi godoc
// @Summary Get transaksi
// @Description Get transaksi
// @Tags exporter
// @Accept  json
// @Produce  json
// @Param request body request.RekapRequest false "..."
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string][]string
// @Router /transaksi-all [post]
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

// Get transaksi by unit id godoc
// @Summary Get transaksi by unit id
// @Description Get transaksi by unit id
// @Tags exporter
// @Accept  json
// @Produce  json
// @Param request body request.RekapRequest false "..."
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string][]string
// @Router /pelanggan [post]
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
