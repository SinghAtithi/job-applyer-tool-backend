package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api/handler"
	"example.com/internal/dto"
	"example.com/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates a test database using SQLite
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	
	return db
}

// SetupRouter sets up a test router
func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	userHandler := handler.NewUserHandler(db)
	
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.PATCH("/:id", userHandler.PatchUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
	
	return router
}

func TestGetAllUsers_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create test users
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Password: "password123"},
		{Name: "Jane Smith", Email: "jane@example.com", Password: "password456"},
	}
	
	for _, user := range users {
		db.Create(&user)
	}
	
	// Test GET /users
	req, _ := http.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetAllUsers_WithPagination(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create test users
	for i := 1; i <= 5; i++ {
		user := models.User{
			Name:     "User",
			Email:    "user",
			Password: "password",
		}
		user.Email = "user" + string(rune(i)) + "@example.com"
		db.Create(&user)
	}
	
	// Test with pagination
	req, _ := http.NewRequest("GET", "/users/?page=1&page_size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password123"}
	db.Create(&user)
	
	// Test GET /users/:id
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetUser_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Test GET /users/999 (non-existent)
	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestCreateUser_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Test POST /users
	userJSON := `{"name": "John Doe", "email": "john@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestCreateUser_ValidationError(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Test POST /users with missing name
	userJSON := `{"email": "john@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestUpdateUser_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password123"}
	db.Create(&user)
	
	// Test PUT /users/1
	userJSON := `{"name": "John Updated", "email": "john.updated@example.com", "password": "newpassword"}`
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestPatchUser_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password123"}
	db.Create(&user)
	
	// Test PATCH /users/1
	patchJSON := `{"name": "John Patched"}`
	req, _ := http.NewRequest("PATCH", "/users/1", bytes.NewBufferString(patchJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestDeleteUser_Success(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupRouter(db)
	
	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password123"}
	db.Create(&user)
	
	// Test DELETE /users/1
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}
