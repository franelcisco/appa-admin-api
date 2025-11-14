package domains

import (
	"context"

	"appa_admin_api/internal/models"
	dbModels "appa_admin_api/pkg/db/models"
)

// LoginService defines the interface for the login service
type LoginService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, uid string) error
}

type OrdersService interface {
	GetManualOrders(ctx context.Context) ([]dbModels.ManualOrder, error)
	UpdateManualOrder(ctx context.Context, req models.UpdateOrderRequest, id int) error
	GetPaymentMethods(ctx context.Context) ([]dbModels.PaymentMethod, error)
}
