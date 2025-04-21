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
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func setupLoginTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db.InitTestDB()
	db.DB = db.TestDB
	auth.RegisterAuthRoutes(r)

	return r
}

func createTestUser(email, password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.DB.Create(&models.User{
		ID:           uuid.New().String(),
		Name:         "Tester",
		Email:        email,
		PasswordHash: string(hash),
	})
}

func TestLoginHandler(t *testing.T) {
	router := setupLoginTestRouter()

	email := "login@example.com"
	password := "test123"
	createTestUser(email, password)

	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
