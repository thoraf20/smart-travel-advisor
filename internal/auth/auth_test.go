package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/thoraf20/smart-travel-advisor/internal/auth"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db.InitTestDB() // 💡 Use SQLite in-memory
	auth.RegisterAuthRoutes(r)

	// point your handlers at db.TestDB instead of db.DB if needed
	db.DB = db.TestDB // 🔁 override global reference

	return r
}

func TestRegisterHandler(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "secure123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}
