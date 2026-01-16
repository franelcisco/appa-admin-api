package shopify

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
)

// Repository defines methods to interact with Shopify API
type Repository interface {
	GetCustomerMetafield(ctx context.Context, customerID string, ns, key string) (*Metafield, error)
	MarkOrderAsPaid(ctx context.Context, gid string) error
}

// Repository is a Shopify API repository
type repository struct {
	gql    *GraphQLClient
	Logger *zap.Logger
}

// NewRepository creates a new Shopify API repository
func NewRepository(
	shopDomain, apiVersion, adminToken string, logger *zap.Logger,
) Repository {
	return &repository{
		gql:    NewGraphQLClient(shopDomain, apiVersion, adminToken, logger),
		Logger: logger,
	}
}

func (r *repository) GetCustomerMetafield(
	ctx context.Context, customerID string, ns, key string,
) (*Metafield, error) {
	gid := GID(customerKind, customerID)
	var resp GetUserByIDResponse
	if err := r.gql.Do(ctx, getCustomerPartnerID, map[string]any{"id": gid, "ns": ns, "key": key}, &resp); err != nil {
		return nil, err
	}
	if resp.Customer == nil || resp.Customer.Metafield == nil {
		return nil, nil
	}

	return resp.Customer.Metafield, nil
}

// MarkOrderAsPaid marks an order as paid
func (r *repository) MarkOrderAsPaid(ctx context.Context, gid string) error {
	if !strings.Contains(gid, OrderKind) {
		gid = GID(OrderKind, gid)
	}

	vars := map[string]any{
		"id": gid,
	}
	var resp MarkOrderAsPaidResponse
	if err := r.gql.Do(ctx, markOrderAsPaid, vars, &resp); err != nil {
		r.Logger.Error(err.Error(), zap.Any("vars", vars))
		return err
	}

	if resp.UserErrors != nil {
		r.Logger.Error("failed to mark order as paid", zap.Any("errors", resp.UserErrors))
		return errors.New("failed to mark order as paid")
	}

	return nil
}
