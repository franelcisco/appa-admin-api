package services

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"appa_admin_api/internal/domains"
	"appa_admin_api/internal/models"
	dbModels "appa_admin_api/pkg/db/models"
	"appa_admin_api/pkg/firebase"
)

type loginService struct {
	Logger   *zap.Logger
	firebase firebase.Repository
	db       *gorm.DB
}

// NewLoginService creates a new login service
func NewLoginService(
	logger *zap.Logger,
	firebase firebase.Repository,
	db *gorm.DB,
) domains.LoginService {
	return &loginService{
		Logger:   logger,
		firebase: firebase,
		db:       db,
	}
}

func (l *loginService) Login(
	ctx context.Context,
	req *models.LoginRequest,
) (*models.LoginResponse, error) {
	firebaseUser, err := l.firebase.GetFirebaseUser(ctx, req.TokeID)
	if err != nil {
		return nil, err
	}

	var dbUser dbModels.User
	err = l.db.Where("email = ?", firebaseUser.Email).First(&dbUser).Error
	if err != nil {
		l.Logger.Error(err.Error(), zap.Any("email", firebaseUser.Email))

		errFirebase := l.firebase.DeleteUser(ctx, firebaseUser.UID)
		if errFirebase != nil {
			l.Logger.Error(errFirebase.Error(), zap.Any("uid", firebaseUser.UID))
		}

		return nil, err
	}

	cookie, err := l.firebase.CreateSessionCookie(ctx, req.TokeID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Email:   dbUser.Email,
		IsAdmin: dbUser.IsAdmin,
		Token:   cookie,
	}, nil
}

// Logout handles user logout requests
func (l *loginService) Logout(ctx context.Context, uid string) error {
	err := l.firebase.RevokeSessionCookie(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}
