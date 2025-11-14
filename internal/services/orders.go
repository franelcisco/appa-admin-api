package services

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"context"

	"appa_admin_api/internal/domains"
	"appa_admin_api/internal/models"
	"appa_admin_api/pkg/db"
	dbModels "appa_admin_api/pkg/db/models"
	"appa_admin_api/pkg/shopify"
)

type ordersService struct {
	db          *gorm.DB
	shopifyRepo shopify.Repository
	logger      *zap.Logger
}

func NewOrdersService(
	db *gorm.DB,
	shopifyRepo shopify.Repository,
	logger *zap.Logger,
) domains.OrdersService {
	return &ordersService{
		db:          db,
		shopifyRepo: shopifyRepo,
		logger:      logger,
	}
}

// GetManualOrders retrieves all manual orders from the database
func (s *ordersService) GetManualOrders(ctx context.Context) ([]dbModels.ManualOrder, error) {

	var orders []dbModels.ManualOrder
	if err := s.db.WithContext(ctx).Find(&orders).Error; err != nil {
		s.logger.Error("failed to retrieve manual orders", zap.Error(err))
		return nil, err
	}

	return orders, nil
}

// UpdateManualOrder updates the validation status of a manual order
func (s *ordersService) UpdateManualOrder(
	ctx context.Context,
	req models.UpdateOrderRequest,
	id int,
) error {
	var (
		errDB error
		tx    = s.db.Begin()
	)
	defer db.DBRollback(tx, &errDB)

	result := tx.Model(&dbModels.ManualOrder{}).WithContext(ctx).
		Where("id = ?", id).
		Updates(map[string]any{
			"validate_status":   req.ValidateStatus,
			"amount":            req.Amount,
			"logistic_validate": req.LogisticValidate,
			"requires_change":   req.RequiresChange,
			"payment_method_id": req.PaymentMethodID,
		})
	if result.Error != nil {
		s.logger.Error("failed to update manual order status", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		s.logger.Warn("no manual order found with the given ID", zap.Int("ID", id))
		return gorm.ErrRecordNotFound
	}

	if req.ValidateStatus == "COMPLETED" {
		errDB := s.shopifyRepo.MarkOrderAsPaid(ctx, req.OrderID)
		if errDB != nil {
			s.logger.Error("failed to mark order as paid in Shopify", zap.Error(errDB))
			return errDB
		}
	}

	return nil
}

// GetPaymentMethods retrieves all payment methods from the database
func (s *ordersService) GetPaymentMethods(ctx context.Context) ([]dbModels.PaymentMethod, error) {
	var methods []dbModels.PaymentMethod
	if err := s.db.WithContext(ctx).Find(&methods).Error; err != nil {
		s.logger.Error("failed to retrieve payment methods", zap.Error(err))
		return nil, err
	}
	return methods, nil
}
