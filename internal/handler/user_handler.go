package handler

import (
	"errors"
	"net/http"

	"github.com/Gergenus/bookingService/internal/dto"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const AccessTokenDuration = 60 * 24 * 60 * 60
const RefreshTokenDuration = 60 * 24 * 60 * 60

type UserHandler struct {
	srv service.UserServiceInterface
}

func NewUserHandler(srv service.UserServiceInterface) UserHandler {
	return UserHandler{srv: srv}
}

func (u *UserHandler) Register(c echo.Context) error {
	var userDTO dto.UserDTO

	err := c.Bind(&userDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	uid, err := u.srv.CreateUser(c.Request().Context(), userDTO.Username, userDTO.Role, userDTO.Email, userDTO.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uid": uid.String(),
	})
}

func (u *UserHandler) Login(c echo.Context) error {
	var loginReq dto.LoginDTO

	err := c.Bind(&loginReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	ip := c.RealIP()
	AccessToken, RefreshToken, err := u.srv.Login(c.Request().Context(), loginReq.Email, loginReq.Password, c.Request().UserAgent(), ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	err = setCookie(c, "AccessToken", AccessToken, AccessTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	err = setCookie(c, "RefreshToken", RefreshToken, RefreshTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"AccessToken":  AccessToken,
		"RefreshToken": RefreshToken,
	})
}

func setCookie(c echo.Context, key, value string, duration int) error {
	cookie := http.Cookie{
		Name:     key,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   duration,
		Value:    value,
		Path:     "/api/v1",
	}
	c.SetCookie(&cookie)
	return nil
}

func (u *UserHandler) Refresh(c echo.Context) error {
	oldRefresh, err := c.Cookie("RefreshToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	refreshUUID, err := uuid.Parse(oldRefresh.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	expiredAccesToken, err := c.Cookie("AccessToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	newRefresh, newAccess, err := u.srv.RefreshToken(c.Request().Context(), refreshUUID, c.Request().UserAgent(), c.RealIP(), expiredAccesToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshSession) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "invalid refresh session",
			})
		}
		if errors.Is(err, service.ErrTokenExpired) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "token expired",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	setCookie(c, "RefreshToken", newRefresh.String(), RefreshTokenDuration)
	setCookie(c, "AccessToken", newAccess, AccessTokenDuration)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"AccessToken":  newAccess,
		"RefreshToken": newRefresh,
	})
}
