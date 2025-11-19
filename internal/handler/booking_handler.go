package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService service.BookingServiceInterface
}

func NewBookingHandler(bookingService service.BookingServiceInterface) BookingHandler {
	return BookingHandler{bookingService: bookingService}
}

func (b *BookingHandler) ScientistBookings(c echo.Context) error {
	uid := c.Get("uuid")
	uuid, ok := uid.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	booking, err := b.bookingService.ScientistBookings(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, booking)
}

func (b *BookingHandler) Createbooking(c echo.Context) error {
	var booking models.Booking
	err := c.Bind(&booking)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid payload",
		})
	}
	uid, ok := c.Get("uuid").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "uuid not found",
		})
	}
	booking.UserId = uuid.MustParse(uid)
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
	uuid, ok := c.Get("uuid").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "uuid not found",
		})
	}
	booking, err := b.bookingService.Booking(c.Request().Context(), bookIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	if booking.UserId.String() != uuid {
		return c.JSON(http.StatusForbidden, map[string]any{
			"error": "invalid owner",
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
