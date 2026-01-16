package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"appa_admin_api/pkg/firebase"
)

// AuthMiddleware provides middleware for authentication
type AuthMiddleware struct {
	firebase firebase.Repository
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(firebase firebase.Repository) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		firebase: firebase,
	}, nil
}

// Auth is a middleware to protect endpoints that require authentication
func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		cookie, err := c.Cookie("SESSION")
		if err != nil {
			if strings.Contains(path, "logout") {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "logged out successfully"})
			} else {
				fmt.Println("Error retrieving cookie:", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			}
			return
		}

		decode, err := a.firebase.VerifySessionCookie(c.Request.Context(), cookie)
		if err != nil {
			if strings.Contains(path, "logout") {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "logged out successfully"})
			} else {
				fmt.Println("Error verifying session cookie:", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			}
			return
		}

		c.Set("UID", decode.UID)

		c.Next()
	}
}
