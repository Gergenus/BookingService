package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/repository"
	"github.com/Gergenus/bookingService/pkg/hash"
	"github.com/Gergenus/bookingService/pkg/jwtpkg"
	"github.com/Gergenus/bookingService/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrPasswordMismatch      = errors.New("password mismatch")
	ErrTokenExpired          = errors.New("token expired")
	ErrInvalidRefreshSession = errors.New("invalid refresh session")
)

type UserService struct {
	userRepo   repository.UserRepositoryInterface
	log        *slog.Logger
	jwtTkn     jwtpkg.TokenService
	RefreshTTL time.Duration
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, role, email, password string) (*uuid.UUID, error)
	Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error)
	RefreshToken(ctx context.Context, oldRefresh uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error)
}

func NewUserService(userRepo repository.UserRepositoryInterface, log *slog.Logger, jwtTkn jwtpkg.TokenService, RefreshTTL time.Duration) *UserService {
	return &UserService{userRepo: userRepo, log: log, jwtTkn: jwtTkn, RefreshTTL: RefreshTTL}
}

func (u *UserService) CreateUser(ctx context.Context, username, role, email, password string) (*uuid.UUID, error) {
	const op = "service.CreateUser"
	u.log.With(slog.String("op", op))
	u.log.Info("Creating user", slog.String("email", email))

	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		slog.Error("failed to generate hashpassword", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user := models.User{
		Username:       username,
		Role:           role,
		Email:          email,
		HashedPassword: hashPassword,
	}

	uid, err := u.userRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			u.log.Error("user already exists error", slog.String("email", user.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		u.log.Error("creating user error", slog.String("email", user.Email), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.log.Info("created user", slog.String("email", user.Email))
	return uid, nil
}

func (u *UserService) Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
	const op = "service.Login"
	u.log.With(slog.String("op", op))
	u.log.Info("logging the user", slog.String("email", email))

	user, err := u.userRepo.UserByEmail(ctx, email)
	if err != nil {
		u.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if !hash.CheckPassword(user.HashedPassword, password) {
		u.log.Warn("passwords mismatch")
		return "", "", fmt.Errorf("%s: %w", op, ErrPasswordMismatch)
	}
	AccessToken, err := u.jwtTkn.GenerateAccessToken(*user)
	if err != nil {
		u.log.Error("failed to create access token", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	RefreshToken := uuid.New()

	fingerprint := utils.CreateFingerprint(ip, userAgent)

	err = u.userRepo.CreateJWTSession(ctx, user.UUID, RefreshToken, fingerprint, ip, time.Now().Add(u.RefreshTTL).Unix(), u.RefreshTTL)
	if err != nil {
		u.log.Error("failed to create jwt session", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessToken, RefreshToken.String(), nil
}

func (u *UserService) RefreshToken(ctx context.Context, oldRefresh uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error) {
	const op = "service.RefreshToken"
	log := u.log.With(slog.String("op", op))
	log.Info("refreshing token")
	session, err := u.userRepo.RefreshSession(ctx, oldRefresh)
	if err != nil {
		log.Error("getting refresh token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	if int64(session.ExpiresIn) < time.Now().Unix() {
		return nil, "", ErrTokenExpired
	}
	fingerprint := utils.CreateFingerprint(ip, userAgent)

	if session.Fingerprint != fingerprint {
		return nil, "", ErrInvalidRefreshSession
	}
	newRefresh := uuid.New()
	err = u.userRepo.CreateJWTSession(ctx, uuid.MustParse(session.UUID), newRefresh, session.Fingerprint,
		session.IP, time.Now().Add(u.RefreshTTL).Unix(), u.RefreshTTL)
	if err != nil {
		log.Error("creating jwt session error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.jwtTkn.RegenerateToken(oldAccessToken)
	if err != nil {
		log.Error("creating jwt token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	return &newRefresh, token, nil
}
