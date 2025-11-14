package routes

import (
	"github.com/gin-gonic/gin"

	"appa_admin_api/internal/handlers"
	"appa_admin_api/pkg/middleware"
)

// LoginRoutes defines the routes for the login service
type LoginRoutes struct {
	handler *handlers.LoginHandler
	auth    *middleware.AuthMiddleware
}

// NewLoginRoutes creates a new instance of LoginRoutes
func NewLoginRoutes(
	handler *handlers.LoginHandler,
	auth *middleware.AuthMiddleware) *LoginRoutes {
	return &LoginRoutes{
		handler: handler,
		auth:    auth,
	}
}

// SetRouter sets up the routes for the login service
func (r *LoginRoutes) SetRouter(router *gin.Engine) {
	router.POST("/login", r.handler.Login)
	router.POST("/login/verify", r.auth.Auth(), r.handler.VerifySession)
	router.POST("/logout", r.handler.Logout)
}
