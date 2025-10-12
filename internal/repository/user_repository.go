package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoSessionFound    = errors.New("no session found")
)

type UserRepository struct {
	db      db.PostgresDB
	redisDB *redis.Client
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*uuid.UUID, error)
	User(ctx context.Context, uuid uuid.UUID) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateJWTSession(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error
	RefreshSession(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error)
}

func NewUserRepository(db db.PostgresDB, redisDB *redis.Client) *UserRepository {
	return &UserRepository{db: db, redisDB: redisDB}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *models.User) (*uuid.UUID, error) {
	const op = "user_repository.CreateUser"

	uuid := uuid.New()
	_, err := u.db.DB.Exec(ctx, "INSERT INTO users (uid, username, role, email, hashed_password) VALUES($1, $2, $3, $4, $5)", uuid.String(),
		user.Username, user.Role, user.Email, user.HashedPassword)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &uuid, nil
}

func (u *UserRepository) User(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	const op = "user_repository.User"
	var user models.User
	err := u.db.DB.QueryRow(ctx, "SELECT * FROM users WHERE uid = $1", uuid.String()).Scan(&user.UUID, &user.Username,
		&user.Role, &user.Email, &user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (u *UserRepository) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.UserByEmail"
	var user models.User
	err := u.db.DB.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(&user.UUID, &user.Username,
		&user.Role, &user.Email, &user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (u *UserRepository) CreateJWTSession(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error {
	const op = "repository.CreateJWTSession"
	values := map[string]any{
		"uuid":          uuid.String(),
		"refresh_token": refreshToken.String(),
		"fingerprint":   fingerprint,
		"ip":            ip,
		"expires_in":    expiresIn,
	}
	err := u.redisDB.HSet(ctx, refreshToken.String(), values).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.redisDB.Expire(ctx, refreshToken.String(), RefreshTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserRepository) RefreshSession(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error) {
	const op = "repository.RefreshSession"
	session, err := u.redisDB.HGetAll(ctx, oldRefresh.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = u.redisDB.Del(ctx, oldRefresh.String()).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exp, err := strconv.Atoi(session["expires_in"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resSession := models.RefreshSession{
		UUID:         session["uuid"],
		RefreshToken: session["refresh_token"],
		Fingerprint:  session["fingerprint"],
		IP:           session["ip"],
		ExpiresIn:    exp,
	}
	return &resSession, nil
}
