package handlers

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"idea-collision-engine-api/internal/auth"
	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/internal/models"
)

type AuthHandler struct {
	db         *database.PostgresDB
	redis      *database.RedisClient
	jwtService *auth.JWTService
	validator  *validator.Validate
}

func NewAuthHandler(db *database.PostgresDB, redis *database.RedisClient, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		db:         db,
		redis:      redis,
		jwtService: jwtService,
		validator:  validator.New(),
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    400,
		})
	}
	
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_failed",
			Message: err.Error(),
			Code:    400,
		})
	}
	
	// Check if user already exists
	existingUser, err := h.db.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to check user existence",
			Code:    500,
		})
	}
	
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:   "user_exists",
			Message: "User with this email already exists",
			Code:    409,
		})
	}
	
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "hash_failed",
			Message: "Failed to hash password",
			Code:    500,
		})
	}
	
	// Create user
	user := &models.User{
		ID:               uuid.New(),
		Email:            req.Email,
		PasswordHash:     hashedPassword,
		SubscriptionTier: models.TierFree,
		Interests:        req.Interests,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	if err := h.db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "user_creation_failed",
			Message: "Failed to create user",
			Code:    500,
		})
	}
	
	// Generate token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate token",
			Code:    500,
		})
	}
	
	// Remove password hash from response
	user.PasswordHash = ""
	
	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    400,
		})
	}
	
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_failed",
			Message: err.Error(),
			Code:    400,
		})
	}
	
	// Get user by email
	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
				Code:    401,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve user",
			Code:    500,
		})
	}
	
	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid email or password",
			Code:    401,
		})
	}
	
	// Generate token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate token",
			Code:    500,
		})
	}
	
	// Remove password hash from response
	user.PasswordHash = ""
	
	return c.JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
			Code:    404,
		})
	}
	
	// Remove password hash from response
	user.PasswordHash = ""
	
	return c.JSON(user)
}

// UpdateProfile updates user profile information
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	
	type UpdateRequest struct {
		Interests []string `json:"interests"`
	}
	
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    400,
		})
	}
	
	// Get current user
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
			Code:    404,
		})
	}
	
	// Update interests
	user.Interests = req.Interests
	user.UpdatedAt = time.Now()
	
	// Note: This would typically use an UpdateUser method
	// For now, we'll return the updated user without persisting
	// You'd need to implement UpdateUser in database layer
	
	user.PasswordHash = ""
	return c.JSON(user)
}