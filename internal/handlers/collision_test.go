package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"idea-collision-engine-api/internal/collision"
	"idea-collision-engine-api/internal/models"
)

// Mock database
type MockPostgresDB struct {
	mock.Mock
}

func (m *MockPostgresDB) GetCollisionDomains(tier string) ([]models.CollisionDomain, error) {
	args := m.Called(tier)
	return args.Get(0).([]models.CollisionDomain), args.Error(1)
}

func (m *MockPostgresDB) CreateCollisionSession(session *models.CollisionSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockPostgresDB) GetUserCollisionHistory(userID uuid.UUID, limit int) ([]models.CollisionSession, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]models.CollisionSession), args.Error(1)
}

func (m *MockPostgresDB) RateCollision(sessionID, userID uuid.UUID, rating int, notes *string) error {
	args := m.Called(sessionID, userID, rating, notes)
	return args.Error(0)
}

func (m *MockPostgresDB) GetUserUsage(userID uuid.UUID) (*models.UserUsage, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserUsage), args.Error(1)
}

func (m *MockPostgresDB) IncrementUserUsage(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Mock Redis
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) GetCachedCollisionDomains(tier string) ([]models.CollisionDomain, error) {
	args := m.Called(tier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CollisionDomain), args.Error(1)
}

func (m *MockRedisClient) CacheCollisionDomains(tier string, domains []models.CollisionDomain, expiration time.Duration) error {
	args := m.Called(tier, domains, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) InvalidateUserUsage(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRedisClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// Mock AI Service
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) EnhanceCollisionResult(result *models.CollisionResult, input models.CollisionInput, domain models.CollisionDomain) error {
	args := m.Called(result, input, domain)
	return args.Error(0)
}

func (m *MockAIService) CheckConnection() error {
	args := m.Called()
	return args.Error(0)
}

type CollisionHandlerTestSuite struct {
	suite.Suite
	app     *fiber.App
	handler *CollisionHandler
	mockDB  *MockPostgresDB
	mockRedis *MockRedisClient
	mockAI  *MockAIService
}

func (suite *CollisionHandlerTestSuite) SetupTest() {
	suite.mockDB = &MockPostgresDB{}
	suite.mockRedis = &MockRedisClient{}
	suite.mockAI = &MockAIService{}
	
	suite.handler = &CollisionHandler{
		db:        suite.mockDB,
		redis:     suite.mockRedis,
		aiService: suite.mockAI,
	}
	
	// Create test domains
	domains := []models.CollisionDomain{
		{
			ID:          uuid.New().String(),
			Name:        "Biomimicry",
			Category:    "Nature",
			Description: "How nature solves problems",
			Keywords:    []string{"evolution", "adaptation"},
			Intensity:   []string{"gentle", "moderate"},
			Tier:        "basic",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Jazz Improvisation",
			Category:    "Music",
			Description: "Spontaneous creation",
			Keywords:    []string{"improvisation", "creativity"},
			Intensity:   []string{"moderate", "radical"},
			Tier:        "basic",
		},
	}
	
	suite.handler.engine = collision.NewCollisionEngine(domains)
	
	// Setup Fiber app
	suite.app = fiber.New()
	suite.setupRoutes()
}

func (suite *CollisionHandlerTestSuite) setupRoutes() {
	// Add middleware to simulate authenticated user
	suite.app.Use(func(c *fiber.Ctx) error {
		userID := uuid.New()
		c.Locals("user_id", userID)
		c.Locals("subscription_tier", models.TierFree)
		return c.Next()
	})
	
	suite.app.Post("/collisions/generate", suite.handler.GenerateCollision)
	suite.app.Get("/collisions/history", suite.handler.GetCollisionHistory)
	suite.app.Put("/collisions/:id/rate", suite.handler.RateCollision)
	suite.app.Get("/collisions/usage", suite.handler.GetUsageStatus)
	suite.app.Get("/collisions/health", suite.handler.HealthCheck)
	suite.app.Get("/domains/basic", suite.handler.GetBasicDomains)
}

func (suite *CollisionHandlerTestSuite) TestGenerateCollision() {
	// Setup mocks
	suite.mockDB.On("CreateCollisionSession", mock.AnythingOfType("*models.CollisionSession")).Return(nil)
	suite.mockDB.On("IncrementUserUsage", mock.AnythingOfType("uuid.UUID")).Return(nil)
	suite.mockRedis.On("InvalidateUserUsage", mock.AnythingOfType("string")).Return(nil)
	
	// Prepare request
	input := models.CollisionInput{
		UserInterests:      []string{"technology", "design"},
		CurrentProject:     "mobile app",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	jsonData, _ := json.Marshal(input)
	req := bytes.NewReader(jsonData)
	
	// Make request
	resp, err := suite.app.Test(suite.createRequest("POST", "/collisions/generate", req))
	
	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	// Parse response
	var result models.CollisionResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	
	assert.NotEmpty(suite.T(), result.ID)
	assert.NotEmpty(suite.T(), result.PrimaryDomain)
	assert.NotEmpty(suite.T(), result.CollisionDomain)
	assert.NotEmpty(suite.T(), result.Connection)
	assert.Greater(suite.T(), len(result.SparkQuestions), 0)
	assert.GreaterOrEqual(suite.T(), result.QualityScore, 0.0)
	
	suite.mockDB.AssertExpectations(suite.T())
	suite.mockRedis.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestGenerateCollisionValidationError() {
	// Invalid input - missing required fields
	invalidInput := models.CollisionInput{
		UserInterests: []string{}, // Empty interests
		// Missing other required fields
	}
	
	jsonData, _ := json.Marshal(invalidInput)
	req := bytes.NewReader(jsonData)
	
	resp, err := suite.app.Test(suite.createRequest("POST", "/collisions/generate", req))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	
	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "validation_failed", errorResp.Error)
}

func (suite *CollisionHandlerTestSuite) TestGetCollisionHistory() {
	userID := uuid.New()
	sessions := []models.CollisionSession{
		{
			ID:     userID,
			UserID: userID,
			CollisionResult: models.CollisionResult{
				ID:              uuid.New().String(),
				PrimaryDomain:   "Technology",
				CollisionDomain: "Jazz Improvisation",
				Connection:      "Test connection",
				QualityScore:    85.5,
			},
			CreatedAt: time.Now(),
		},
	}
	
	suite.mockDB.On("GetUserCollisionHistory", mock.AnythingOfType("uuid.UUID"), 20).Return(sessions, nil)
	
	resp, err := suite.app.Test(suite.createRequest("GET", "/collisions/history", nil))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result []models.CollisionSession
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Jazz Improvisation", result[0].CollisionResult.CollisionDomain)
	
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestRateCollision() {
	sessionID := uuid.New()
	
	suite.mockDB.On("RateCollision", sessionID, mock.AnythingOfType("uuid.UUID"), 5, mock.AnythingOfType("*string")).Return(nil)
	
	rateRequest := map[string]interface{}{
		"rating": 5,
		"notes":  "Great collision!",
	}
	
	jsonData, _ := json.Marshal(rateRequest)
	req := bytes.NewReader(jsonData)
	
	resp, err := suite.app.Test(suite.createRequest("PUT", "/collisions/"+sessionID.String()+"/rate", req))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Rating saved successfully", result["message"])
	
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestGetUsageStatus() {
	userID := uuid.New()
	usage := &models.UserUsage{
		ID:             uuid.New(),
		UserID:         userID,
		CollisionCount: 3,
		ResetDate:      time.Now(),
	}
	
	suite.mockDB.On("GetUserUsage", mock.AnythingOfType("uuid.UUID")).Return(usage, nil)
	
	resp, err := suite.app.Test(suite.createRequest("GET", "/collisions/usage", nil))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), models.TierFree, result["tier"])
	assert.Equal(suite.T(), float64(3), result["collisions_used"])
	assert.Equal(suite.T(), float64(5), result["collisions_limit"])
	assert.Equal(suite.T(), float64(2), result["collisions_remaining"])
	assert.Equal(suite.T(), false, result["unlimited"])
	
	suite.mockDB.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestGetBasicDomains() {
	domains := []models.CollisionDomain{
		{
			Name:        "Biomimicry",
			Category:    "Nature",
			Description: "How nature solves problems",
			Tier:        "basic",
		},
	}
	
	// Test cache miss - should query database
	suite.mockRedis.On("GetCachedCollisionDomains", "basic").Return(nil, nil)
	suite.mockDB.On("GetCollisionDomains", "basic").Return(domains, nil)
	suite.mockRedis.On("CacheCollisionDomains", "basic", domains, 30*time.Minute).Return(nil)
	
	resp, err := suite.app.Test(suite.createRequest("GET", "/domains/basic", nil))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result []models.CollisionDomain
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Biomimicry", result[0].Name)
	
	suite.mockDB.AssertExpectations(suite.T())
	suite.mockRedis.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestHealthCheck() {
	// Setup mocks for healthy services
	suite.mockRedis.On("Ping").Return(nil)
	suite.mockAI.On("CheckConnection").Return(nil)
	
	resp, err := suite.app.Test(suite.createRequest("GET", "/collisions/health", nil))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "collision-engine", result["service"])
	assert.Equal(suite.T(), "healthy", result["status"])
	assert.Equal(suite.T(), "connected", result["database"])
	assert.Equal(suite.T(), "connected", result["cache"])
	assert.Equal(suite.T(), "connected", result["ai_service"])
	assert.Equal(suite.T(), "ready", result["collision_engine"])
	
	suite.mockRedis.AssertExpectations(suite.T())
	suite.mockAI.AssertExpectations(suite.T())
}

func (suite *CollisionHandlerTestSuite) TestHealthCheckDegraded() {
	// Setup mocks for degraded services
	suite.mockRedis.On("Ping").Return(assert.AnError)
	suite.mockAI.On("CheckConnection").Return(assert.AnError)
	
	resp, err := suite.app.Test(suite.createRequest("GET", "/collisions/health", nil))
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "degraded", result["status"])
	assert.Equal(suite.T(), "unavailable", result["cache"])
	assert.Equal(suite.T(), "unavailable", result["ai_service"])
}

// Helper method to create HTTP requests
func (suite *CollisionHandlerTestSuite) createRequest(method, path string, body *bytes.Reader) *http.Request {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	return req
}

func TestCollisionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CollisionHandlerTestSuite))
}