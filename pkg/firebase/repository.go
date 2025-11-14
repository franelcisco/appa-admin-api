package firebase

import (
	"context"
	"errors"
	"time"

	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
)

// Repository defines methods to interact with Firebase API
type Repository interface {
	GetFirebaseUser(ctx context.Context, tokenID string) (*User, error)
	CreateSessionCookie(ctx context.Context, tokenID string) (string, error)
	VerifySessionCookie(ctx context.Context, cookie string) (*User, error)
	RevokeSessionCookie(ctx context.Context, uid string) error
	DeleteUser(ctx context.Context, uid string) error
}

type repository struct {
	app    *firebase.App
	logger *zap.Logger
}

// NewRepository creates a new Firebase repository
func NewRepository(
	logger *zap.Logger, serviceAccountID, projectID string,
) (Repository, error) {
	conf := &firebase.Config{
		ServiceAccountID: serviceAccountID,
		ProjectID:        projectID,
	}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return &repository{
		app:    app,
		logger: logger,
	}, nil
}

func (r *repository) GetFirebaseUser(
	ctx context.Context,
	tokenID string,
) (*User, error) {
	client, err := r.app.Auth(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, tokenID)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	user := &User{
		UID: token.UID,
	}
	if email, ok := token.Claims["email"]; ok && email.(string) != "" {
		user.Email = email.(string)
	} else {
		r.logger.Error("not email in user")
		return nil, errors.New("not email in user")
	}

	return user, nil
}

func (r *repository) CreateSessionCookie(
	ctx context.Context,
	tokenID string,
) (string, error) {
	client, err := r.app.Auth(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return "", err
	}

	expires := 12 * time.Hour

	cookie, err := client.SessionCookie(ctx, tokenID, expires)
	if err != nil {
		r.logger.Error(err.Error())
		return "", err
	}

	return cookie, nil
}

func (r *repository) VerifySessionCookie(
	ctx context.Context,
	cookie string,
) (*User, error) {
	client, err := r.app.Auth(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	token, err := client.VerifySessionCookieAndCheckRevoked(ctx, cookie)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	user := &User{
		UID: token.UID,
	}
	if email, ok := token.Claims["email"]; ok && email.(string) != "" {
		user.Email = email.(string)
	} else {
		return nil, errors.New("not email in user")
	}

	return user, nil
}

// RevokeSessionCookie revokes all refresh tokens for a specified user identified by their UID.
func (r *repository) RevokeSessionCookie(
	ctx context.Context, uid string,
) error {
	client, err := r.app.Auth(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	if err := client.RevokeRefreshTokens(ctx, uid); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}

// DeleteUser deletes a user by their UID
func (r *repository) DeleteUser(ctx context.Context, uid string) error {
	client, err := r.app.Auth(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	if err := client.DeleteUser(ctx, uid); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}
