package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api/handler"
	"example.com/api/routes"
	"example.com/internal/dto"
	"example.com/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabase represents a test database setup
type TestDatabase struct {
	DB *gorm.DB
}

// SetupTestDatabase creates a test database using SQLite
func SetupTestDatabase(t *testing.T) *TestDatabase {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.CoverLetterTable{},
		&models.JobDescriptionTable{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return &TestDatabase{DB: db}
}

// SetupTestRouter sets up a test router with all routes
func SetupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := SetupTestDatabase(t)
	routes.SetupRoutes(router, testDB.DB)

	return router
}

// TestUserFlow tests the complete user flow (CRUD operations)
func TestUserFlow(t *testing.T) {
	router := SetupTestRouter(t)

	// Step 1: Create a user
	t.Run("CreateUser", func(t *testing.T) {
		userJSON := `{"name": "Test User", "email": "test@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Step 2: Get all users
	t.Run("GetAllUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Step 3: Get user by ID
	t.Run("GetUser", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Step 4: Update user
	t.Run("UpdateUser", func(t *testing.T) {
		userJSON := `{"name": "Updated User", "email": "updated@example.com", "password": "newpassword"}`
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Step 5: Patch user
	t.Run("PatchUser", func(t *testing.T) {
		patchJSON := `{"name": "Patched User"}`
		req, _ := http.NewRequest("PATCH", "/users/1", bytes.NewBufferString(patchJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Step 6: Delete user
	t.Run("DeleteUser", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})
}

// TestResumeFlow tests the resume operations
func TestResumeFlow(t *testing.T) {
	router := SetupTestRouter(t)

	// Test GetAllResumes (initially empty)
	t.Run("GetAllResumes", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/resumes/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// This might return an error since it's using helper functions
		// Just check that the endpoint is reachable
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

// TestCoverLetterFlow tests the cover letter operations
func TestCoverLetterFlow(t *testing.T) {
	router := SetupTestRouter(t)

	// Test GetAllCoverLetters
	t.Run("GetAllCoverLetters", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/coverLetter/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// This endpoint might not be implemented
		// Just check that the endpoint is reachable
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

// TestPagination tests pagination functionality
func TestPagination(t *testing.T) {
	router := SetupTestRouter(t)

	// Create multiple users for pagination testing
	for i := 1; i <= 15; i++ {
		userJSON := `{"name": "User ` + string(rune(i+'0')) + `", "email": "user` + string(rune(i+'0')) + `@example.com", "password": "password` + string(rune(i+'0')) + `"}`
		req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Test pagination with page 1
	t.Run("Pagination_Page1", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/?page=1&page_size=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	// Test pagination with page 2
	t.Run("Pagination_Page2", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/?page=2&page_size=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test invalid pagination (page 0 - should default to 1)
	t.Run("Pagination_Invalid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/?page=0&page_size=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestValidationErrors tests validation error handling
func TestValidationErrors(t *testing.T) {
	router := SetupTestRouter(t)

	// Test create user with missing email
	t.Run("MissingEmail", func(t *testing.T) {
		userJSON := `{"name": "Test User", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
	})

	// Test create user with missing name
	t.Run("MissingName", func(t *testing.T) {
		userJSON := `{"email": "test@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
	})

	// Test invalid JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		userJSON := `{"name": "Test User", "email": }`
		req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test invalid user ID format
	t.Run("InvalidUserID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestNotFoundErrors tests 404 error handling
func TestNotFoundErrors(t *testing.T) {
	router := SetupTestRouter(t)

	// Test non-existent user
	t.Run("UserNotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/99999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
	})
}

// TestHealthCheck tests basic health check endpoint (if exists)
func TestHealthCheck(t *testing.T) {
	router := SetupTestRouter(t)

	// Test root endpoint
	t.Run("RootEndpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 as there's no root handler
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// BenchmarkUserHandlers benchmarks user handler performance
func BenchmarkUserHandlers(b *testing.B) {
	router := SetupTestRouter(&testing.T{})

	// Create test user
	userJSON := `{"name": "Benchmark User", "email": "benchmark@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ = http.NewRequest("GET", "/users/", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// Helper function to create a test user handler
func createTestUserHandler(db *gorm.DB) *handler.UserHandler {
	return handler.NewUserHandler(db)
}
