package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the JWT token
		userToken := c.Get("user")
		if userToken == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing JWT token")
		}

		token, ok := userToken.(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT token")
		}

		claims, ok := token.Claims.(*domain.JwtClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT claims")
		}

		// Parse subject (user ID)
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID in token: "+err.Error())
		}

		// Store userID and claims in context
		c.Set("userID", userID)
		c.Set("claims", claims)
		return next(c)
	}
}

func RequirePermission(required domain.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get claims from context
			claimsVal := c.Get("claims")
			if claimsVal == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing JWT claims")
			}

			claims, ok := claimsVal.(*domain.JwtClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims type")
			}

			// Check if required permission is present
			for _, perm := range strings.Split(claims.Scope, " ") {
				if domain.Permission(perm) == domain.PermAll || domain.Permission(perm) == required {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("missing required permission: %s", required))
		}
	}
}
