package controllers

import (
	"fmt"
	"net/http"
	"yoga-guru/internal/config"
	"yoga-guru/internal/middleware"
	"yoga-guru/internal/models"
	"yoga-guru/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler provides methods for authentication.
type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

// RegisterRequest defines the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=student instructor"`
	Gender   string `json:"gender" binding:"omitempty,oneof=male female"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password. Role defaults to 'student'.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string "message: User registered successfully"
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 409 {object} map[string]string "error: User with this email or username already exists"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userRole := models.Student // Default role is student
	if req.Role != "" {
		// In a real application, only an admin should be able to set roles other than student.
		// For now, we allow it for testing, but typically this would be restricted.
		switch models.UserRole(req.Role) {
		case models.Admin, models.Instructor, models.Student:
			userRole = models.UserRole(req.Role)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
			return
		}
	}

	user := models.User{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         userRole,
		Profile: models.Profile{
			Name:   req.Name,
			Gender: models.UserGender(req.Gender),
		},
	}

	fmt.Print(user.Profile.Name, user.Profile.Gender)

	// Check if user already exists
	var existingUser models.User
	if h.DB.Where("phone = ?", user.Phone).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginRequest defines the request body for user login.
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user with email and password, returning a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]string "token: JWT_TOKEN"
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Invalid credentials"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, refresh, err := middleware.GenerateJWT(user.ID.String(), user.Role, h.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh": refresh, "role": user.Role})
}

// RefreshTokenRequest defines the request body for refreshing the token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshToken godoc
// @Summary Refresh an access token
// @Description Refresh a JWT access token using a valid refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string "token: JWT_TOKEN, refresh: REFRESH_TOKEN"
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Invalid refresh token"
// @Failure 404 {object} map[string]string "error: User not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := middleware.ValidateToken(req.RefreshToken, h.Cfg)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", claims.Subject).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, newRefresh, err := middleware.GenerateJWT(user.ID.String(), user.Role, h.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh": newRefresh, "role": user.Role})
}
