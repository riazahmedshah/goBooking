package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/server"

	echojwt "github.com/labstack/echo-jwt/v4"
)

const (
	UserIDKey = "userID"
	RoleKey   = "userRole"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(server *server.Server) *AuthMiddleware {
	return &AuthMiddleware{server: server}
}

func (auth *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(auth.server.Config.JWT.SecretKey),
		TokenLookup: "cookie:access_token",
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
	
		return jwtMiddleware(func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok || token == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			userID, _ := claims["userId"].(string)
			role, _ := claims["role"].(string)

			c.Set(UserIDKey, userID)
			c.Set(RoleKey, role)

			return next(c)
		})
	}
}

func (auth *AuthMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok || token == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			userRole, ok := claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing role in token claims")
			}

			for _, role := range roles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "you do not have the required permissions to access this resource")
		}
	}
}
