package controllers

import (
	"net/http"
	"strconv"
	"yoga-guru/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler provides methods for user management.
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type UserProfileResponse struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatarURL"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
}

// GetCurrentUserProfile godoc
// @Summary Get current user's profile
// @Description Retrieve the profile details of the authenticated user.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserProfileResponse
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: User not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUserProfile(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDAny.(string)

	var user models.User
	if err := h.DB.Model(&models.User{ID: uuid.MustParse(userID)}).Preload("Profile").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		Name:      user.Profile.Name,
		AvatarURL: user.Profile.AvatarURL,
		Phone:     user.Phone,
		Gender:    string(user.Profile.Gender),
	})
}

// UpdateUserRoleRequest defines the request body for updating a user's role.
type UpdateUserRoleRequest struct {
	Role string
}

// UpdateUserRole godoc
// @Summary Update a user's role (Admin only)
// @Description Allows an admin to update the role of any user.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param role body UpdateUserRoleRequest true "New role for the user"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: User not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userIDParam, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	targetUserID := uint(userIDParam)

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, targetUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	newRole := models.UserRole(req.Role)
	switch newRole {
	case models.Admin, models.Instructor, models.Student:
		// Valid role
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	user.Role = newRole
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	user.PasswordHash = "" // Don't expose password hash
	c.JSON(http.StatusOK, user)
}
