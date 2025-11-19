package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Gergenus/bookingService/internal/dto"
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
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "photo.jpg")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "form file error",
		})
	}
	file, err := image.Open()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "form file error",
		})
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "copying error",
		})
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	req, err := http.NewRequest("POST", "http://host.docker.internal:5000/describe", body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	if resp.StatusCode == 400 {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	desc := dto.MLResponse{}
	err = json.Unmarshal(data, &desc)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	eq.Description = desc.Description
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

func (e *EquipmentHandler) SignedImageURL(c echo.Context) error {
	imagePath := c.Param("image")
	obj, err := e.srv.SignURL(c.Request().Context(), imagePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
	}
	c.Response().Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(c.Response().Writer, obj)
	return err
}
