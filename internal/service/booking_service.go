package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/repository"
)

var (
	ErrIntervalInterception = errors.New("interval interception")
)

type BookingService struct {
	bookingRepo repository.BookingRepositoryInterface
	log         *slog.Logger
}

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, booking models.Booking) (int, error)
	Bookings(ctx context.Context, equipmentId int) ([]models.Booking, error)
	DeleteBooking(ctx context.Context, bookingId int) error
	Booking(ctx context.Context, bookingId int) (*models.Booking, error)
	ScientistBookings(ctx context.Context, uid string) ([]models.Booking, error)
}

func NewBookingService(bookingRepo repository.BookingRepositoryInterface, log *slog.Logger) BookingService {
	return BookingService{bookingRepo: bookingRepo, log: log}
}

func (b *BookingService) ScientistBookings(ctx context.Context, uid string) ([]models.Booking, error) {
	booking, err := b.bookingRepo.ScientistBookings(ctx, uid)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (b *BookingService) CreateBooking(ctx context.Context, booking models.Booking) (int, error) {
	const op = "booking_service.CreateBooking"
	log := b.log.With(slog.String("op", op))
	log.Info("creating booking", slog.Int("equipment_id", booking.EquipmentId), slog.String("user_id", booking.UserId.String()))
	id, err := b.bookingRepo.CreateBooking(ctx, booking)
	if err != nil {
		if errors.Is(err, repository.ErrIntervalInterception) {
			return 0, fmt.Errorf("%s: %w", op, ErrIntervalInterception)
		}
		log.Error("creating booking error", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (b *BookingService) Bookings(ctx context.Context, equipmentId int) ([]models.Booking, error) {
	const op = "booking_service.Bookings"
	log := b.log.With(slog.String("op", op))
	log.Info("getting bookings", slog.Int("equipment_id", equipmentId))
	bookings, err := b.bookingRepo.Bookings(ctx, equipmentId)
	if err != nil {
		log.Error("getting bookings error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return bookings, nil
}

func (b *BookingService) DeleteBooking(ctx context.Context, bookingId int) error {
	const op = "booking_service.DeleteBooking"
	log := b.log.With(slog.String("op", op))
	log.Info("deleting bookings", slog.Int("booking_id", bookingId))
	err := b.bookingRepo.DeleteBooking(ctx, bookingId)
	if err != nil {
		log.Error("deleting bookings error", slog.Int("booking_id", bookingId))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (b *BookingService) Booking(ctx context.Context, bookingId int) (*models.Booking, error) {
	const op = "booking_service.Booking"
	log := b.log.With(slog.String("op", op))
	log.Info("getting booking", slog.Int("booking_id", bookingId))
	booking, err := b.bookingRepo.Booking(ctx, bookingId)
	if err != nil {
		log.Error("getting bookings error", slog.Int("booking_id", bookingId))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return booking, nil
}
