package middleware

import (
	"net/http"
	"strings"

	"github/Rubncal04/youtube-premium/auth"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip authentication for public routes
			if c.Path() == "/" || c.Path() == "/login" || c.Path() == "/refresh" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := auth.ValidateToken(tokenString, secretKey)
			if err != nil {
				if err == auth.ErrExpiredToken {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "token has expired",
					})
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
			}

			// Set the user ID in the context
			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}
