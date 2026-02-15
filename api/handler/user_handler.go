package handler

import (
	"log"
	"net/http"
	"strconv"

	"example.com/internal/dto"
	"example.com/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	gorm "gorm.io/gorm"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetAllUsers retrieves all users with pagination
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query users
	var users []models.User
	var total int64

	// Get total count
	if err := h.db.Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("failed to count users: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to count users",
			"An internal error occurred",
		))
		return
	}

	// Get paginated results (deterministic ordering)
	if err := h.db.Order("id ASC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		log.Printf("failed to fetch users: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch users",
			"An internal error occurred",
		))
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, dto.SuccessResponse(
		dto.NewPaginatedResponse(users, page, pageSize, total),
		"Users retrieved successfully",
	))
}

// GetUser retrieves a single user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"User ID is required",
			"",
		))
		return
	}

	// parse and validate id to avoid accidental WHERE clause injection
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid user ID format",
			"",
		))
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"NOT_FOUND",
				"User not found",
				"",
			))
			return
		}
		log.Printf("failed to fetch user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch user",
			"An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User retrieved successfully"))
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	// bind to input DTO to avoid mass-assignment
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR",
			"Failed to process request",
			"An internal error occurred",
		))
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to create user",
			"An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(user, "User created successfully"))
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"User ID is required",
			"",
		))
		return
	}

	// Parse ID
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid user ID format",
			"",
		))
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"NOT_FOUND",
				"User not found",
				"",
			))
			return
		}
		log.Printf("failed to fetch user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch user",
			"An internal error occurred",
		))
		return
	}

	// Bind updated data
	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Update fields
	user.Name = updateData.Name
	user.Email = updateData.Email
	if updateData.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to hash password: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
				"INTERNAL_ERROR",
				"Failed to process request",
				"An internal error occurred",
			))
			return
		}
		user.Password = string(hashed)
	}

	if err := h.db.Save(&user).Error; err != nil {
		log.Printf("failed to save updated user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to update user",
			"An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User updated successfully"))
}

// PatchUser partially updates a user
func (h *UserHandler) PatchUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"User ID is required",
			"",
		))
		return
	}

	// Parse ID
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid user ID format",
			"",
		))
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"NOT_FOUND",
				"User not found",
				"",
			))
			return
		}
		log.Printf("failed to fetch user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch user",
			"An internal error occurred",
		))
		return
	}

	// Bind patch data
	var patchData map[string]interface{}
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	// Apply patches
	if name, ok := patchData["name"].(string); ok {
		user.Name = name
	}
	if email, ok := patchData["email"].(string); ok {
		user.Email = email
	}
	if password, ok := patchData["password"].(string); ok {
		// Hash patched password before storing
		hashed, herr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if herr != nil {
			log.Printf("failed to hash patched password (id=%d): %v", userID, herr)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
				"INTERNAL_ERROR",
				"Failed to process request",
				"An internal error occurred",
			))
			return
		}
		user.Password = string(hashed)
	}

	if err := h.db.Save(&user).Error; err != nil {
		log.Printf("failed to save patched user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to update user",
			"An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User patched successfully"))
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"User ID is required",
			"",
		))
		return
	}

	// Parse ID
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid user ID format",
			"",
		))
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"NOT_FOUND",
				"User not found",
				"",
			))
			return
		}
		log.Printf("failed to fetch user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch user",
			"An internal error occurred",
		))
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		log.Printf("failed to delete user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to delete user",
			"An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "User deleted successfully"))
}
