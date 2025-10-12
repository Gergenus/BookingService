package jwtpkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrClaimsFailed         = errors.New("claims failed")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type JWTService struct {
	Secret    string
	AccessTTL time.Duration
}

type TokenService interface {
	GenerateAccessToken(user models.User) (string, error)
	RegenerateToken(oldToken string) (string, error)
	ParseToken(token string) (*UserClaims, error)
}

func NewUserJWTpkg(Secret string, AccessTTL time.Duration) JWTService {
	return JWTService{Secret: Secret, AccessTTL: AccessTTL}
}

func (j JWTService) GenerateAccessToken(user models.User) (string, error) {
	const op = "jwtpkg.GenerateAccessToken"

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTTL)),
		},
		UUID:     user.UUID.String(),
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	AccessTokenString, err := AccessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessTokenString, nil
}

func (u JWTService) RegenerateToken(oldToken string) (string, error) {
	const op = "jwtpkg.RegenerateToken"
	var claims UserClaims
	token, err := jwt.ParseWithClaims(oldToken, &claims, func(t *jwt.Token) (interface{}, error) { return []byte(u.Secret), nil })
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	parsedUUID, err := uuid.Parse(claims.UUID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return "", ErrClaimsFailed
	}

	user := models.User{
		UUID:     parsedUUID,
		Username: claims.Username,
		Role:     claims.Role,
		Email:    claims.Email,
	}
	return u.GenerateAccessToken(user)
}

// returns role, uuid and error
func (j JWTService) ParseToken(token string) (*UserClaims, error) {
	const op = "jwtpkg.ParseToken"
	var claims UserClaims

	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !tkn.Valid {
		return nil, ErrClaimsFailed
	}

	return &claims, nil
}
