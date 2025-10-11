package handler

import (
	"net/http"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/labstack/echo/v4"
)

type EquipmentHandler struct {
	srv service.EquipmentService
}

func NewEquipmentHandler(srv service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{srv: srv}
}

func (e *EquipmentHandler) CreateEquipment(c echo.Context) error {
	var eq models.Equipment
	err := c.Bind(&eq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid request",
		})
	}
	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid request" + err.Error(),
		})
	}
	id, err := e.srv.CreateEquipment(c.Request().Context(), eq, image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}
