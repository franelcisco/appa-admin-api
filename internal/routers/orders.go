package routes

import (
	"github.com/gin-gonic/gin"

	"appa_admin_api/internal/handlers"
	"appa_admin_api/pkg/middleware"
)

type OrdersRoutes struct {
	handler *handlers.OrdersHandler
	auth    *middleware.AuthMiddleware
}

func NewOrdersRoutes(
	handler *handlers.OrdersHandler,
	auth *middleware.AuthMiddleware,
) *OrdersRoutes {
	return &OrdersRoutes{
		handler: handler,
		auth:    auth,
	}
}

func (r *OrdersRoutes) SetRouter(router *gin.Engine) {
	router.GET("/bcv-tasa", r.handler.GetBCVTasa)
	authorized := router.Group("/orders")
	authorized.Use(r.auth.Auth())
	{
		authorized.GET("/manual", r.handler.GetManualOrders)
		authorized.PUT("/manual/:id", r.handler.UpdateManualOrder)
		authorized.GET("/payment-methods", r.handler.GetPaymentMethods)
	}
}
