package middleware

import "github.com/riazahmedshah/go-booking/internal/server"

type Middlewares struct {
	Auth *AuthMiddleware
}

func NewMiddleware(server *server.Server) *Middlewares {
	return &Middlewares{
		Auth: NewAuthMiddleware(server),
	}
}
