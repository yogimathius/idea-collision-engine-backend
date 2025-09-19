package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/internal/models"
)

type RateLimitConfig struct {
	WindowSeconds int
	MaxRequests   int
	SkipPremium   bool
}

// RateLimitMiddleware implements rate limiting using Redis
func RateLimitMiddleware(redis *database.RedisClient, config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by auth middleware)
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			// If no user ID, use IP address for rate limiting
			userID, _ := uuid.Parse(c.IP())
			if userID == uuid.Nil {
				return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
					Error:   "rate_limit_exceeded",
					Message: "Unable to identify user for rate limiting",
					Code:    429,
				})
			}
		}

		userIDStr := userID.String()

		// Skip rate limiting for premium users if configured
		if config.SkipPremium {
			tier := GetSubscriptionTierFromContext(c)
			if tier == models.TierPro || tier == models.TierTeam {
				return c.Next()
			}
		}

		// Check rate limit
		allowed, err := redis.CheckRateLimit(userIDStr, config.WindowSeconds, config.MaxRequests)
		if err != nil {
			// Log error but don't block request if Redis is down
			fmt.Printf("Rate limit check failed: %v\n", err)
			return c.Next()
		}

		if !allowed {
			// Get rate limit status for headers
			remaining, resetTime, _ := redis.GetRateLimitStatus(userIDStr, config.WindowSeconds, config.MaxRequests)
			
			c.Set("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
			c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Set("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(resetTime).Unix())))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: fmt.Sprintf("Rate limit exceeded. Try again in %v seconds", int(resetTime.Seconds())),
				Code:    429,
			})
		}

		// Set rate limit headers
		remaining, resetTime, _ := redis.GetRateLimitStatus(userIDStr, config.WindowSeconds, config.MaxRequests)
		c.Set("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining-1)) // -1 for current request
		c.Set("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(resetTime).Unix())))

		return c.Next()
	}
}

// UsageLimitMiddleware checks freemium collision limits
func UsageLimitMiddleware(db *database.PostgresDB, redis *database.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			return err
		}

		tier := GetSubscriptionTierFromContext(c)
		
		// Skip usage limits for premium users
		if tier == models.TierPro || tier == models.TierTeam {
			return c.Next()
		}

		// Check cached usage first
		userIDStr := userID.String()
		usage, err := redis.GetCachedUserUsage(userIDStr)
		if err != nil {
			// Fallback to database
			usage, err = db.GetUserUsage(userID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
					Error:   "usage_check_failed",
					Message: "Unable to check usage limits",
					Code:    500,
				})
			}
			
			// Cache the result
			redis.CacheUserUsage(userIDStr, usage, 5*time.Minute)
		}

		// Check if user has exceeded weekly limit
		limit := models.UsageLimits[tier]
		if limit > 0 && usage.CollisionCount >= limit {
			return c.Status(fiber.StatusPaymentRequired).JSON(models.ErrorResponse{
				Error:   "usage_limit_exceeded",
				Message: fmt.Sprintf("Weekly limit of %d collisions exceeded. Upgrade to Pro for unlimited access.", limit),
				Code:    402,
			})
		}

		// Store usage in context for handlers to increment
		c.Locals("user_usage", usage)

		return c.Next()
	}
}