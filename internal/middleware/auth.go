package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"idea-collision-engine-api/internal/auth"
	"idea-collision-engine-api/internal/models"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header required",
				Code:    401,
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid authorization header format",
				Code:    401,
			})
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid token",
				Code:    401,
			})
		}

		// Store user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("subscription_tier", claims.SubscriptionTier)

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Next()
		}

		// Store user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("subscription_tier", claims.SubscriptionTier)

		return c.Next()
	}
}

// GetUserIDFromContext extracts user ID from Fiber context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	return id, nil
}

// GetSubscriptionTierFromContext extracts subscription tier from context
func GetSubscriptionTierFromContext(c *fiber.Ctx) string {
	tier := c.Locals("subscription_tier")
	if tier == nil {
		return models.TierFree
	}

	tierStr, ok := tier.(string)
	if !ok {
		return models.TierFree
	}

	return tierStr
}

// RequirePremium middleware requires pro or team subscription
func RequirePremium() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tier := GetSubscriptionTierFromContext(c)
		
		if tier != models.TierPro && tier != models.TierTeam {
			return c.Status(fiber.StatusPaymentRequired).JSON(models.ErrorResponse{
				Error:   "premium_required",
				Message: "This feature requires a premium subscription",
				Code:    402,
			})
		}

		return c.Next()
	}
}