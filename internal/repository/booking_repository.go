package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/pkg/db"
)

var (
	ErrIntervalInterception = errors.New("interval interception")
)

type PostgresBookingRepository struct {
	db db.PostgresDB
}

type BookingRepositoryInterface interface {
	CreateBooking(ctx context.Context, booking models.Booking) (int, error)
	Bookings(ctx context.Context, equipmentId int) ([]models.Booking, error)
	DeleteBooking(ctx context.Context, bookingId int) error
	Booking(ctx context.Context, bookingId int) (*models.Booking, error)
	ScientistBookings(ctx context.Context, uid string) ([]models.Booking, error)
}

func NewPostgresBookingRepository(db db.PostgresDB) PostgresBookingRepository {
	return PostgresBookingRepository{db: db}
}

func (p *PostgresBookingRepository) ScientistBookings(ctx context.Context, uid string) ([]models.Booking, error) {
	var data []models.Booking
	row, err := p.db.DB.Query(ctx, "SELECT * FROM booking WHERE user_id = $1", uid)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var booking models.Booking
		err := row.Scan(&booking.Id, &booking.EquipmentId, &booking.UserId, &booking.StartTime, &booking.EndTime)
		if err != nil {
			return nil, err
		}
		data = append(data, booking)
	}
	return data, nil
}

func (p *PostgresBookingRepository) checkInterceptions(ctx context.Context, startTime, endTime time.Time, equipmentId int) (bool, error) {
	const op = "booking_repository.checkInterceprions"
	var count int
	err := p.db.DB.QueryRow(ctx, "SELECT Count(*) FROM booking WHERE equipment_id = $1 AND ((start_time <= $3 AND end_time >= $3)"+
		"OR ($2 >= start_time AND $2 <= end_time) OR ($2 <= start_time AND $3 >= end_time) OR ($2 >= start_time AND $3 <= end_time))", equipmentId, startTime, endTime).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return count == 0, nil
}

func (p *PostgresBookingRepository) CreateBooking(ctx context.Context, booking models.Booking) (int, error) {
	const op = "booking_repository.CreateBooking"
	ok, err := p.checkInterceptions(ctx, booking.StartTime, booking.EndTime, booking.EquipmentId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, ErrIntervalInterception)
	}
	var id int
	err = p.db.DB.QueryRow(ctx, "INSERT INTO booking (equipment_id, user_id, start_time, end_time) VALUES($1, $2, $3, $4) RETURNING id",
		booking.EquipmentId, booking.UserId, booking.StartTime, booking.EndTime).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *PostgresBookingRepository) Bookings(ctx context.Context, equipmentId int) ([]models.Booking, error) {
	const op = "booking_repository.Bookings"
	var bookings []models.Booking
	rows, err := p.db.DB.Query(ctx, "SELECT * FROM booking WHERE equipment_id = $1", equipmentId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.Id, &booking.EquipmentId, &booking.UserId, &booking.StartTime, &booking.EndTime)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bookings = append(bookings, booking)
	}
	rows.Close()
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return bookings, nil
}

func (p *PostgresBookingRepository) DeleteBooking(ctx context.Context, bookingId int) error {
	const op = "booking_repository.DeleteBooking"
	_, err := p.db.DB.Exec(ctx, "DELETE FROM booking WHERE id = $1", bookingId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *PostgresBookingRepository) Booking(ctx context.Context, bookingId int) (*models.Booking, error) {
	const op = "booking_repository.Booking"
	var booking models.Booking
	err := p.db.DB.QueryRow(ctx, "SELECT * FROM booking WHERE id = $1", bookingId).Scan(&booking.Id, &booking.EquipmentId,
		&booking.UserId, &booking.StartTime, &booking.EndTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &booking, nil
}
