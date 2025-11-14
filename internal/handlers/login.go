package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"appa_admin_api/internal/domains"
	"appa_admin_api/internal/models"
)

// LoginHandler handles user login requests
type LoginHandler struct {
	loginService domains.LoginService
	debug        bool
}

// NewLoginHandler creates a new instance of LoginHandler
func NewLoginHandler(loginService domains.LoginService, debug string) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
		debug:        debug == "1",
	}
}

// Login handles user login requests
func (h *LoginHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.loginService.Login(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("Login error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Set cookie attributes based on debug mode
	var secure bool
	var sameSite http.SameSite
	if h.debug {
		// HTTP: not secure, lax same-site
		secure = false
		sameSite = http.SameSiteLaxMode
	} else {
		// HTTPS: secure, none same-site
		fmt.Println("Request is over HTTPS")
		secure = true
		sameSite = http.SameSiteNoneMode
	}
	// create session Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "SESSION",
		Value:    resp.Token,
		Path:     "/",
		Domain:   "",
		MaxAge:   int(12 * time.Hour / time.Second),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	c.JSON(http.StatusOK, resp)
}

// VerifySession verifies if the current session is valid.
func (h *LoginHandler) VerifySession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "session valid"})
}

// Logout handles user logout requests
func (h *LoginHandler) Logout(c *gin.Context) {
	uid, exists := c.Get("UID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.loginService.Logout(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	// Clear the session cookie
	c.SetCookie(
		"SESSION",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
