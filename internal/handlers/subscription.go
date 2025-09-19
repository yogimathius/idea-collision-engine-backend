package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"

	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/internal/middleware"
	"idea-collision-engine-api/internal/models"
)

type SubscriptionHandler struct {
	db    *database.PostgresDB
	redis *database.RedisClient
}

func NewSubscriptionHandler(db *database.PostgresDB, redis *database.RedisClient, stripeKey string) *SubscriptionHandler {
	stripe.Key = stripeKey
	
	return &SubscriptionHandler{
		db:    db,
		redis: redis,
	}
}

// Stripe price IDs (these would be configured in Stripe dashboard)
const (
	ProMonthlyPriceID  = "price_pro_monthly"  // Replace with actual Stripe price ID
	TeamMonthlyPriceID = "price_team_monthly" // Replace with actual Stripe price ID
)

// CreateCheckoutSession creates a Stripe checkout session for subscription
func (h *SubscriptionHandler) CreateCheckoutSession(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	type CheckoutRequest struct {
		PriceID     string `json:"price_id" validate:"required"`
		SuccessURL  string `json:"success_url" validate:"required"`
		CancelURL   string `json:"cancel_url" validate:"required"`
	}
	
	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    400,
		})
	}
	
	// Validate price ID
	if req.PriceID != ProMonthlyPriceID && req.PriceID != TeamMonthlyPriceID {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_price_id",
			Message: "Invalid price ID",
			Code:    400,
		})
	}
	
	// Get user details
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
			Code:    404,
		})
	}
	
	// Create or get Stripe customer
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}
	
	stripeCustomer, err := customer.New(customerParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "customer_creation_failed",
			Message: "Failed to create Stripe customer",
			Code:    500,
		})
	}
	
	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(stripeCustomer.ID),
		SuccessURL: stripe.String(req.SuccessURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(req.CancelURL),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(req.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}
	
	session, err := session.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "session_creation_failed",
			Message: "Failed to create checkout session",
			Code:    500,
		})
	}
	
	return c.JSON(fiber.Map{
		"checkout_url": session.URL,
		"session_id":   session.ID,
	})
}

// GetSubscriptionStatus returns the current subscription status
func (h *SubscriptionHandler) GetSubscriptionStatus(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
			Code:    404,
		})
	}
	
	response := fiber.Map{
		"tier":        user.SubscriptionTier,
		"status":      "active",
		"expires_at":  nil,
		"cancel_at":   nil,
		"is_trial":    false,
		"features": h.getTierFeatures(user.SubscriptionTier),
	}
	
	// For premium users, we'd typically store and retrieve Stripe subscription details
	// This is a simplified version
	if user.SubscriptionTier == models.TierPro || user.SubscriptionTier == models.TierTeam {
		response["billing_cycle"] = "monthly"
		response["next_billing_date"] = time.Now().AddDate(0, 1, 0)
	}
	
	return c.JSON(response)
}

// CancelSubscription cancels the user's subscription
func (h *SubscriptionHandler) CancelSubscription(c *fiber.Ctx) error {
	_, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}
	
	// In a real implementation, you'd:
	// 1. Find the Stripe subscription ID for this user
	// 2. Cancel the subscription via Stripe API
	// 3. Update the user's tier in the database
	
	// For now, return a placeholder response
	return c.JSON(fiber.Map{
		"message": "Subscription cancellation initiated",
		"status":  "cancelled",
		"effective_date": time.Now().AddDate(0, 1, 0), // End of current billing period
	})
}

// WebhookHandler handles Stripe webhooks
func (h *SubscriptionHandler) WebhookHandler(c *fiber.Ctx) error {
	// Get the webhook signature
	sig := c.Get("Stripe-Signature")
	if sig == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "missing_signature",
			Message: "Missing Stripe signature",
			Code:    400,
		})
	}
	
	// Get the request body
	body := c.Body()
	
	// In a real implementation, you'd:
	// 1. Verify the webhook signature
	// 2. Parse the webhook event
	// 3. Handle different event types (customer.subscription.created, etc.)
	// 4. Update user subscription status in database
	
	// For now, return success
	fmt.Printf("Received Stripe webhook: %s\n", string(body))
	
	return c.JSON(fiber.Map{
		"received": true,
	})
}

// GetPricingPlans returns available pricing plans
func (h *SubscriptionHandler) GetPricingPlans(c *fiber.Ctx) error {
	plans := []fiber.Map{
		{
			"id":          "free",
			"name":        "Free",
			"price":       0,
			"currency":    "usd",
			"interval":    "month",
			"price_id":    nil,
			"features": []string{
				"5 collisions per week",
				"Basic collision domains",
				"Standard quality insights",
			},
			"limits": fiber.Map{
				"collisions_per_week": 5,
				"premium_domains":     false,
				"ai_enhancement":      false,
			},
		},
		{
			"id":          "pro",
			"name":        "Pro",
			"price":       12,
			"currency":    "usd",
			"interval":    "month",
			"price_id":    ProMonthlyPriceID,
			"features": []string{
				"Unlimited collisions",
				"Premium collision domains",
				"AI-enhanced insights",
				"Advanced spark questions",
				"Priority support",
			},
			"limits": fiber.Map{
				"collisions_per_week": -1,
				"premium_domains":     true,
				"ai_enhancement":      true,
			},
		},
		{
			"id":          "team",
			"name":        "Team",
			"price":       39,
			"currency":    "usd",
			"interval":    "month",
			"price_id":    TeamMonthlyPriceID,
			"features": []string{
				"Everything in Pro",
				"Team collaboration features",
				"Custom collision domains",
				"Analytics and insights",
				"Priority support",
			},
			"limits": fiber.Map{
				"collisions_per_week": -1,
				"premium_domains":     true,
				"ai_enhancement":      true,
				"team_features":       true,
			},
		},
	}
	
	return c.JSON(fiber.Map{
		"plans": plans,
	})
}

// getTierFeatures returns features available for a subscription tier
func (h *SubscriptionHandler) getTierFeatures(tier string) []string {
	switch tier {
	case models.TierFree:
		return []string{
			"5 collisions per week",
			"Basic collision domains",
			"Standard quality insights",
		}
	case models.TierPro:
		return []string{
			"Unlimited collisions",
			"Premium collision domains",
			"AI-enhanced insights",
			"Advanced spark questions",
			"Priority support",
		}
	case models.TierTeam:
		return []string{
			"Everything in Pro",
			"Team collaboration features",
			"Custom collision domains",
			"Analytics and insights",
			"Priority support",
		}
	default:
		return []string{}
	}
}

// UpdateSubscriptionTier updates user's subscription tier (called from webhook)
func (h *SubscriptionHandler) UpdateSubscriptionTier(userID uuid.UUID, tier string) error {
	// In a real implementation, you'd update the user's subscription tier in the database
	// For now, this is a placeholder
	fmt.Printf("Updating user %s to tier %s\n", userID, tier)
	
	// Invalidate user cache
	h.redis.InvalidateUserUsage(userID.String())
	
	return nil
}