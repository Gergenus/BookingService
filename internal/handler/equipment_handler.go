package handler

import (
	"net/http"
	"strconv"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/labstack/echo/v4"
)

type EquipmentHandler struct {
	srv service.EquipmentServiceInterface
}

// form-file: image
func NewEquipmentHandler(srv service.EquipmentServiceInterface) *EquipmentHandler {
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
			"error": "invalid request",
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

func (e *EquipmentHandler) DeleteEquipment(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "bad request",
		})
	}
	err = e.srv.DeleteEquipment(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func (e *EquipmentHandler) EquipmentById(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "bad request",
		})
	}
	eq, err := e.srv.Equipment(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, eq)
}

func (e *EquipmentHandler) EquipmentByName(c echo.Context) error {
	name := c.QueryParam("query")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid payload",
		})
	}
	eqs, err := e.srv.EquipmentByName(c.Request().Context(), name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, eqs)
}
