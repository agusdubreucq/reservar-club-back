package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/adubr/reservar-club-back/pkg/jwt"
	"github.com/adubr/reservar-club-back/pkg/oauth2"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService    *services.AuthService
	tokenManager   *jwt.TokenManager
	googleProvider *oauth2.GoogleOAuth2Provider
}

func NewAuthHandler(authService *services.AuthService, tokenManager *jwt.TokenManager, googleProvider *oauth2.GoogleOAuth2Provider) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		tokenManager:   tokenManager,
		googleProvider: googleProvider,
	}
}

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (h *AuthHandler) GetAuthURL(c *gin.Context) {
	state := generateRandomString(32)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL := h.googleProvider.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.googleProvider.ExchangeCode(c.Request.Context(), req.Code)
	if err != nil {
		log.Printf("failed to exchange code: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange code"})
		return
	}

	userInfo, err := h.googleProvider.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		log.Printf("failed to get user info: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get user info"})
		return
	}

	user, err := h.authService.GetOrCreateUserByGoogle(c.Request.Context(), userInfo.Email, userInfo.ID, userInfo.Name)
	if err != nil {
		log.Printf("failed to get or create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	jwtToken, err := h.tokenManager.GenerateToken(user.ID, user.Email, 24*60*60)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:   jwtToken,
		UserID:  user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Message: "login successful",
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
