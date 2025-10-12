package middleware

import (
	"net/http"

	"github.com/Gergenus/bookingService/pkg/jwtpkg"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	JWT jwtpkg.TokenService
}

func NewJWTMiddleware(JWT jwtpkg.TokenService) *JWTMiddleware {
	return &JWTMiddleware{JWT: JWT}
}

func (j *JWTMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("AccessToken")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "getting token error",
			})
		}
		if cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "No auth token",
			})
		}
		claims, err := j.JWT.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}
		c.Set("role", claims.Role)
		c.Set("uuid", claims.UUID)
		return next(c)
	}
}

func (j *JWTMiddleware) AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("role")
		userRoleStr, ok := userRole.(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "user role not found",
			})
		}
		if userRoleStr != "admin" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "admin access required",
			})
		}
		return next(c)
	}
}

func (j *JWTMiddleware) ScientistAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("role")
		userRoleStr, ok := userRole.(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "user role not found",
			})
		}
		if userRoleStr != "scientist" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "admin access required",
			})
		}
		return next(c)
	}
}
