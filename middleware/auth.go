package middleware

import (
	"net/http"
	"strings"

	"github/Rubncal04/youtube-premium/auth"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip authentication for public routes
			if c.Request().URL.Path == "/register" || c.Request().URL.Path == "/login" || c.Request().URL.Path == "/refresh" || c.Request().URL.Path == "/" {
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

			// Convert userID string to ObjectID
			userID, err := primitive.ObjectIDFromHex(claims.UserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid user ID in token",
				})
			}

			// Set the user ID in the context
			c.Set("user_id", userID)

			return next(c)
		}
	}
}
