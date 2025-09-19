package models

import (
	"time"

	"github.com/google/uuid"
)

// CollisionInput represents the user input for collision generation
type CollisionInput struct {
	UserInterests      []string `json:"user_interests" validate:"required,min=1"`
	CurrentProject     string   `json:"current_project" validate:"required"`
	ProjectType        string   `json:"project_type" validate:"required,oneof=product content business research"`
	CollisionIntensity string   `json:"collision_intensity" validate:"required,oneof=gentle moderate radical"`
}

// CollisionResult represents the generated collision output
type CollisionResult struct {
	ID              string    `json:"id" db:"id"`
	PrimaryDomain   string    `json:"primary_domain" db:"primary_domain"`
	CollisionDomain string    `json:"collision_domain" db:"collision_domain"`
	Connection      string    `json:"connection" db:"connection"`
	SparkQuestions  []string  `json:"spark_questions" db:"spark_questions"`
	Examples        []string  `json:"examples" db:"examples"`
	NextSteps       []string  `json:"next_steps" db:"next_steps"`
	QualityScore    float64   `json:"quality_score" db:"quality_score"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
	Rating          *int      `json:"rating,omitempty" db:"rating"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
}

// CollisionDomain represents a curated domain for collision generation
type CollisionDomain struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Category    string   `json:"category" db:"category"`
	Description string   `json:"description" db:"description"`
	Examples    []string `json:"examples" db:"examples"`
	Keywords    []string `json:"keywords" db:"keywords"`
	Intensity   []string `json:"intensity" db:"intensity"`
	Tier        string   `json:"tier" db:"tier"` // basic, premium, custom
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	PasswordHash     string    `json:"-" db:"password_hash"`
	SubscriptionTier string    `json:"subscription_tier" db:"subscription_tier"` // free, pro, team
	Interests        []string  `json:"interests" db:"interests"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// CollisionSession represents a collision generation session
type CollisionSession struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	UserID           uuid.UUID       `json:"user_id" db:"user_id"`
	InputData        CollisionInput  `json:"input_data" db:"input_data"`
	CollisionResult  CollisionResult `json:"collision_result" db:"collision_result"`
	UserRating       *int            `json:"user_rating,omitempty" db:"user_rating"`
	ExplorationNotes *string         `json:"exploration_notes,omitempty" db:"exploration_notes"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// UserUsage represents user usage tracking for freemium limits
type UserUsage struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	CollisionCount int       `json:"collision_count" db:"collision_count"`
	ResetDate      time.Time `json:"reset_date" db:"reset_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	Interests []string `json:"interests,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// SubscriptionTier constants
const (
	TierFree = "free"
	TierPro  = "pro"
	TierTeam = "team"
)

// Usage limits per tier per week
var UsageLimits = map[string]int{
	TierFree: 5,
	TierPro:  -1, // unlimited
	TierTeam: -1, // unlimited
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}