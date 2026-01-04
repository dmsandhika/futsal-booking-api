package handler

import (
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"futsal-booking/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthHandler struct {
	Repo *repository.AdminRepository
}

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var admin *model.Admin
	var err error
	if req.Identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identity required"})
		return
	}

	// Try username first, then email
	admin, err = h.Repo.GetAdminByUsername(req.Identity)
	if err != nil {
		admin, err = h.Repo.GetAdminByEmail(req.Identity)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid identity or password"})
		return
	}

	// Compare hashed password
	if bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid identity or password"})
		return
	}

	token, err := utils.GenerateJWT(int(admin.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.PasswordConfirm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if req.Password != req.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check if username or email already exists
	if _, err := h.Repo.GetAdminByUsername(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}
	if _, err := h.Repo.GetAdminByEmail(req.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	admin := &model.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.Repo.CreateAdmin(admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var adminID uint
	switch v := adminIDVal.(type) {
	case float64:
		adminID = uint(v)
	case int:
		adminID = uint(v)
	}

	admin, err := h.Repo.GetAdminByID(adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       admin.ID,
		"username": admin.Username,
		"email":    admin.Email,
	})
}
