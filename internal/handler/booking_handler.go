package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService service.BookingServiceInterface
}

func NewBookingHandler(bookingService service.BookingServiceInterface) BookingHandler {
	return BookingHandler{bookingService: bookingService}
}

func (b *BookingHandler) Createbooking(c echo.Context) error {
	var booking models.Booking
	err := c.Bind(&booking)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid payload",
		})
	}
	// TODO: Добавить замену uid из middleware
	id, err := b.bookingService.CreateBooking(c.Request().Context(), booking)
	if err != nil {
		if errors.Is(err, service.ErrIntervalInterception) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": "interval interception",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

func (b *BookingHandler) Bookings(c echo.Context) error {
	eqId := c.Param("id")
	eqIdInt, err := strconv.Atoi(eqId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid payload",
		})
	}
	bookings, err := b.bookingService.Bookings(c.Request().Context(), eqIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, bookings)
}

func (b *BookingHandler) DeleteBooking(c echo.Context) error {
	bookId := c.Param("id")
	bookIdInt, err := strconv.Atoi(bookId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid payload",
		})
	}
	err = b.bookingService.DeleteBooking(c.Request().Context(), bookIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}
