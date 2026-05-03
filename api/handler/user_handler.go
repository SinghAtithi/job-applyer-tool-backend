package handler

import (
	"net/http"

	"example.com/internal/dto"
	"example.com/internal/models"
	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// fetchUserByID looks up a user by ID, sending the appropriate error response if not found.
// Returns the user and true on success.
func (h *UserHandler) fetchUserByID(c *gin.Context, userID uint64) (*models.User, bool) {
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", "User not found", ""))
			return nil, false
		}
		logger.Error("failed to fetch user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to fetch user", "An internal error occurred",
		))
		return nil, false
	}
	return &user, true
}

// GetAllUsers retrieves all users with pagination, sorting, and field selection
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var queryParams dto.QueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR", "Invalid query parameters", err.Error(),
		))
		return
	}

	offset := queryParams.GetOffset()
	limit := queryParams.GetLimit()
	page := queryParams.Page
	if page == 0 {
		page = 1
	}
	pageSize := queryParams.PageSize
	if pageSize == 0 {
		pageSize = limit
	}

	var users []models.User
	var total int64

	if err := h.db.Model(&models.User{}).Count(&total).Error; err != nil {
		logger.Error("failed to count users: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to count users", "An internal error occurred",
		))
		return
	}

	query := h.db.Model(&models.User{})

	// Apply sorting
	if sortField, isDesc := queryParams.GetSortFieldAndOrder(); sortField != "" {
		order := "ASC"
		if isDesc {
			order = "DESC"
		}
		query = query.Order(sortField + " " + order)
	} else {
		query = query.Order("id ASC")
	}

	// Apply field selection
	if fields := queryParams.GetFieldsSlice(); fields != nil {
		query = query.Select(fields)
	}

	// Apply pagination
	if offset > 0 || limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		logger.Error("failed to fetch users: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to fetch users", "An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		dto.NewPaginatedResponse(users, page, pageSize, total),
		"Users retrieved successfully",
	))
}

// GetUser retrieves a single user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, ok := parseIDParam(c)
	if !ok {
		return
	}

	user, ok := h.fetchUserByID(c, userID)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User retrieved successfully"))
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR", "Invalid request body", err.Error(),
		))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR", "Failed to process request", "An internal error occurred",
		))
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	if err := h.db.Create(&user).Error; err != nil {
		logger.Error("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to create user", "An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(user, "User created successfully"))
}

// UpdateUser fully replaces an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, ok := parseIDParam(c)
	if !ok {
		return
	}

	user, ok := h.fetchUserByID(c, userID)
	if !ok {
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR", "Invalid request body", err.Error(),
		))
		return
	}

	user.Name = updateData.Name
	user.Email = updateData.Email
	if updateData.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("failed to hash password: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
				"INTERNAL_ERROR", "Failed to process request", "An internal error occurred",
			))
			return
		}
		user.Password = string(hashed)
	}

	if err := h.db.Save(user).Error; err != nil {
		logger.Error("failed to update user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to update user", "An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User updated successfully"))
}

// PatchUser partially updates a user
func (h *UserHandler) PatchUser(c *gin.Context) {
	userID, ok := parseIDParam(c)
	if !ok {
		return
	}

	user, ok := h.fetchUserByID(c, userID)
	if !ok {
		return
	}

	var patchData map[string]interface{}
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"VALIDATION_ERROR", "Invalid request body", err.Error(),
		))
		return
	}

	if name, ok := patchData["name"].(string); ok && name != "" {
		user.Name = name
	}
	if email, ok := patchData["email"].(string); ok && email != "" {
		user.Email = email
	}
	if password, ok := patchData["password"].(string); ok && password != "" {
		hashed, herr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if herr != nil {
			logger.Error("failed to hash patched password (id=%d): %v", userID, herr)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
				"INTERNAL_ERROR", "Failed to process request", "An internal error occurred",
			))
			return
		}
		user.Password = string(hashed)
	}

	if err := h.db.Save(user).Error; err != nil {
		logger.Error("failed to save patched user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to update user", "An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(user, "User patched successfully"))
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, ok := parseIDParam(c)
	if !ok {
		return
	}

	user, ok := h.fetchUserByID(c, userID)
	if !ok {
		return
	}

	if err := h.db.Delete(user).Error; err != nil {
		logger.Error("failed to delete user (id=%d): %v", userID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR", "Failed to delete user", "An internal error occurred",
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "User deleted successfully"))
}
