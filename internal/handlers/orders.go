package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"appa_admin_api/internal/domains"
	"appa_admin_api/internal/models"
	"appa_admin_api/pkg/bcv"
)

type OrdersHandler struct {
	OrdersService domains.OrdersService
	bcvClient     bcv.Client
}

func NewOrdersHandler(ordersService domains.OrdersService, bcvClient bcv.Client) *OrdersHandler {
	return &OrdersHandler{
		OrdersService: ordersService,
		bcvClient:     bcvClient,
	}
}

// GetBCVTasa handles the request to get the BCV exchange rate
func (h *OrdersHandler) GetBCVTasa(c *gin.Context) {
	rate, err := h.bcvClient.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BCVTasaUSDResponse{
		Date: time.Now().Format("2006-01-02"),
		Rate: rate,
	})
}

// GetManualOrders handles the retrieval of manual orders
func (h *OrdersHandler) GetManualOrders(c *gin.Context) {
	orders, err := h.OrdersService.GetManualOrders(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateManualOrder handles updating manual order
func (h *OrdersHandler) UpdateManualOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.OrdersService.UpdateManualOrder(context.Background(), req, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// GetPaymentMethods handles the retrieval of payment methods
func (h *OrdersHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.OrdersService.GetPaymentMethods(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, methods)
}
