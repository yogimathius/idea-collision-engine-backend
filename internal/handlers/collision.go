package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"idea-collision-engine-api/internal/collision"
	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/internal/middleware"
	"idea-collision-engine-api/internal/models"
)

type CollisionHandler struct {
	db         *database.PostgresDB
	redis      *database.RedisClient
	engine     *collision.CollisionEngine
	aiService  *collision.AIService
	validator  *validator.Validate
}

func NewCollisionHandler(db *database.PostgresDB, redis *database.RedisClient, aiService *collision.AIService) *CollisionHandler {
	return &CollisionHandler{
		db:        db,
		redis:     redis,
		aiService: aiService,
		validator: validator.New(),
	}
}

// Initialize loads collision domains and creates the engine
func (h *CollisionHandler) Initialize() error {
	// Load all domains for basic tier (covers all users)
	domains, err := h.db.GetCollisionDomains("premium") // Get all domains
	if err != nil {
		return fmt.Errorf("failed to load collision domains: %w", err)
	}
	
	h.engine = collision.NewCollisionEngine(domains)
	return nil
}

// GenerateCollision creates a new collision for the user
func (h *CollisionHandler) GenerateCollision(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	tier := middleware.GetSubscriptionTierFromContext(c)
	
	var input models.CollisionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    400,
		})
	}
	
	if err := h.validator.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_failed",
			Message: err.Error(),
			Code:    400,
		})
	}
	
	// Generate collision
	result, err := h.engine.GenerateCollision(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "collision_generation_failed",
			Message: "Failed to generate collision",
			Code:    500,
		})
	}
	
	// Enhance with AI for premium users
	if tier == models.TierPro || tier == models.TierTeam {
		domain := h.findDomainByName(result.CollisionDomain)
		if domain != nil {
			if err := h.aiService.EnhanceCollisionResult(result, input, *domain); err != nil {
				// Log error but don't fail the request
				fmt.Printf("AI enhancement failed: %v\n", err)
			}
		}
	}
	
	// Save collision session
	session := &models.CollisionSession{
		ID:              uuid.New(),
		UserID:          userID,
		InputData:       input,
		CollisionResult: *result,
		CreatedAt:       time.Now(),
	}
	
	if err := h.db.CreateCollisionSession(session); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save collision session: %v\n", err)
	}
	
	// Increment usage for free tier users
	if tier == models.TierFree {
		if err := h.db.IncrementUserUsage(userID); err != nil {
			fmt.Printf("Failed to increment usage: %v\n", err)
		}
		
		// Invalidate cache
		h.redis.InvalidateUserUsage(userID.String())
	}
	
	return c.JSON(result)
}

// GetCollisionHistory returns user's collision history
func (h *CollisionHandler) GetCollisionHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	// Parse limit parameter
	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	
	sessions, err := h.db.GetUserCollisionHistory(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve collision history",
			Code:    500,
		})
	}
	
	return c.JSON(sessions)
}

// RateCollision allows users to rate and add notes to their collisions
func (h *CollisionHandler) RateCollision(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	sessionIDStr := c.Params("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_session_id",
			Message: "Invalid session ID",
			Code:    400,
		})
	}
	
	type RateRequest struct {
		Rating int     `json:"rating" validate:"required,min=1,max=5"`
		Notes  *string `json:"notes,omitempty"`
	}
	
	var req RateRequest
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
	
	// Update the rating
	if err := h.db.RateCollision(sessionID, userID, req.Rating, req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "rating_failed",
			Message: "Failed to save rating",
			Code:    500,
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Rating saved successfully",
	})
}

// GetPremiumDomains returns premium domains for Pro/Team users
func (h *CollisionHandler) GetPremiumDomains(c *fiber.Ctx) error {
	tier := middleware.GetSubscriptionTierFromContext(c)
	
	// Check premium access
	if tier != models.TierPro && tier != models.TierTeam {
		return c.Status(fiber.StatusPaymentRequired).JSON(models.ErrorResponse{
			Error:   "premium_required",
			Message: "Premium subscription required to access premium domains",
			Code:    402,
		})
	}
	
	// Try cache first
	cachedDomains, err := h.redis.GetCachedCollisionDomains("premium")
	if err == nil && cachedDomains != nil {
		return c.JSON(cachedDomains)
	}
	
	// Get from database
	domains, err := h.db.GetCollisionDomains("premium")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve premium domains",
			Code:    500,
		})
	}
	
	// Cache the result
	h.redis.CacheCollisionDomains("premium", domains, 30*time.Minute)
	
	return c.JSON(domains)
}

// GetBasicDomains returns basic domains available to all users
func (h *CollisionHandler) GetBasicDomains(c *fiber.Ctx) error {
	// Try cache first
	cachedDomains, err := h.redis.GetCachedCollisionDomains("basic")
	if err == nil && cachedDomains != nil {
		return c.JSON(cachedDomains)
	}
	
	// Get from database
	domains, err := h.db.GetCollisionDomains("basic")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve basic domains",
			Code:    500,
		})
	}
	
	// Cache the result
	h.redis.CacheCollisionDomains("basic", domains, 30*time.Minute)
	
	return c.JSON(domains)
}

// GetUsageStatus returns current usage information for the user
func (h *CollisionHandler) GetUsageStatus(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	tier := middleware.GetSubscriptionTierFromContext(c)
	
	// Premium users have unlimited usage
	if tier == models.TierPro || tier == models.TierTeam {
		return c.JSON(fiber.Map{
			"tier":               tier,
			"collisions_used":    0,
			"collisions_limit":   -1,
			"collisions_remaining": -1,
			"reset_date":         nil,
			"unlimited":          true,
		})
	}
	
	// Get usage for free tier users
	usage, err := h.db.GetUserUsage(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "usage_check_failed",
			Message: "Failed to check usage",
			Code:    500,
		})
	}
	
	limit := models.UsageLimits[tier]
	remaining := limit - usage.CollisionCount
	if remaining < 0 {
		remaining = 0
	}
	
	return c.JSON(fiber.Map{
		"tier":                tier,
		"collisions_used":     usage.CollisionCount,
		"collisions_limit":    limit,
		"collisions_remaining": remaining,
		"reset_date":          usage.ResetDate,
		"unlimited":           false,
	})
}

// findDomainByName helper function to find a domain by name
func (h *CollisionHandler) findDomainByName(name string) *models.CollisionDomain {
	if h.engine == nil {
		return nil
	}
	
	for _, domain := range h.engine.Domains {
		if domain.Name == name {
			return &domain
		}
	}
	
	return nil
}

// HealthCheck endpoint for collision service
func (h *CollisionHandler) HealthCheck(c *fiber.Ctx) error {
	status := fiber.Map{
		"service":   "collision-engine",
		"status":    "healthy",
		"timestamp": time.Now(),
	}
	
	// Check database connection
	if h.db == nil {
		status["database"] = "unavailable"
		status["status"] = "unhealthy"
	} else {
		status["database"] = "connected"
	}
	
	// Check Redis connection
	if err := h.redis.Ping(); err != nil {
		status["cache"] = "unavailable"
		status["status"] = "degraded"
	} else {
		status["cache"] = "connected"
	}
	
	// Check AI service
	if err := h.aiService.CheckConnection(); err != nil {
		status["ai_service"] = "unavailable"
		status["status"] = "degraded"
	} else {
		status["ai_service"] = "connected"
	}
	
	// Check collision engine
	if h.engine == nil {
		status["collision_engine"] = "uninitialized"
		status["status"] = "unhealthy"
	} else {
		status["collision_engine"] = "ready"
	}
	
	statusCode := fiber.StatusOK
	if status["status"] == "unhealthy" {
		statusCode = fiber.StatusServiceUnavailable
	}
	
	return c.Status(statusCode).JSON(status)
}