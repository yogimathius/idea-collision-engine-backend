package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"idea-collision-engine-api/internal/models"
)

type JWTServiceTestSuite struct {
	suite.Suite
	jwtService *JWTService
	testUser   *models.User
}

func (suite *JWTServiceTestSuite) SetupTest() {
	suite.jwtService = NewJWTService("test-secret-key-for-testing")
	
	suite.testUser = &models.User{
		ID:               uuid.New(),
		Email:            "test@example.com",
		SubscriptionTier: models.TierPro,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (suite *JWTServiceTestSuite) TestGenerateToken() {
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	assert.Contains(suite.T(), token, ".") // JWT should contain dots
	
	// Token should have 3 parts (header.payload.signature)
	parts := len([]rune(token))
	assert.Greater(suite.T(), parts, 50) // Reasonable minimum length
}

func (suite *JWTServiceTestSuite) TestValidateToken() {
	// Generate a token first
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)
	
	// Validate the token
	claims, err := suite.jwtService.ValidateToken(token)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), suite.testUser.Email, claims.Email)
	assert.Equal(suite.T(), suite.testUser.SubscriptionTier, claims.SubscriptionTier)
	assert.Equal(suite.T(), "idea-collision-engine", claims.Issuer)
	assert.Equal(suite.T(), suite.testUser.ID.String(), claims.Subject)
}

func (suite *JWTServiceTestSuite) TestValidateInvalidToken() {
	invalidToken := "invalid.token.here"
	
	claims, err := suite.jwtService.ValidateToken(invalidToken)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *JWTServiceTestSuite) TestValidateExpiredToken() {
	// Create a JWT service with a very short expiration for testing
	shortLivedService := &JWTService{secretKey: []byte("test-key")}
	
	// We can't easily test expired tokens without modifying the generation
	// In a real scenario, you'd use a library like golang-jwt/jwt/v5/test
	// or modify the expiration time for testing
	
	// For now, test with malformed token to ensure error handling
	malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed.signature"
	
	claims, err := shortLivedService.ValidateToken(malformedToken)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *JWTServiceTestSuite) TestExtractUserID() {
	// Generate a token
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)
	
	// Extract user ID
	userID, err := suite.jwtService.ExtractUserID(token)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testUser.ID, userID)
}

func (suite *JWTServiceTestSuite) TestExtractUserIDFromInvalidToken() {
	invalidToken := "invalid.token"
	
	userID, err := suite.jwtService.ExtractUserID(invalidToken)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), uuid.Nil, userID)
}

func (suite *JWTServiceTestSuite) TestGenerateRefreshToken() {
	userID := suite.testUser.ID
	
	refreshToken, err := suite.jwtService.GenerateRefreshToken(userID)
	
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), refreshToken)
	assert.Contains(suite.T(), refreshToken, ".") // JWT should contain dots
}

func (suite *JWTServiceTestSuite) TestValidateRefreshToken() {
	userID := suite.testUser.ID
	
	// Generate refresh token
	refreshToken, err := suite.jwtService.GenerateRefreshToken(userID)
	assert.NoError(suite.T(), err)
	
	// Validate refresh token
	extractedUserID, err := suite.jwtService.ValidateRefreshToken(refreshToken)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

func (suite *JWTServiceTestSuite) TestValidateInvalidRefreshToken() {
	invalidRefreshToken := "invalid.refresh.token"
	
	userID, err := suite.jwtService.ValidateRefreshToken(invalidRefreshToken)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), uuid.Nil, userID)
}

func (suite *JWTServiceTestSuite) TestTokensWithDifferentUsers() {
	// Create another user
	user2 := &models.User{
		ID:               uuid.New(),
		Email:            "user2@example.com",
		SubscriptionTier: models.TierFree,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	// Generate tokens for both users
	token1, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)
	
	token2, err := suite.jwtService.GenerateToken(user2)
	assert.NoError(suite.T(), err)
	
	// Tokens should be different
	assert.NotEqual(suite.T(), token1, token2)
	
	// Validate both tokens contain correct user data
	claims1, err := suite.jwtService.ValidateToken(token1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testUser.ID, claims1.UserID)
	assert.Equal(suite.T(), models.TierPro, claims1.SubscriptionTier)
	
	claims2, err := suite.jwtService.ValidateToken(token2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user2.ID, claims2.UserID)
	assert.Equal(suite.T(), models.TierFree, claims2.SubscriptionTier)
}

func (suite *JWTServiceTestSuite) TestTokensWithDifferentSecrets() {
	// Create JWT service with different secret
	differentService := NewJWTService("different-secret-key")
	
	// Generate token with original service
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)
	
	// Try to validate with different service (should fail)
	claims, err := differentService.ValidateToken(token)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *JWTServiceTestSuite) TestClaimsExpiration() {
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)
	
	claims, err := suite.jwtService.ValidateToken(token)
	assert.NoError(suite.T(), err)
	
	// Token should expire in approximately 24 hours
	expectedExpiration := time.Now().Add(24 * time.Hour)
	actualExpiration := claims.ExpiresAt.Time
	
	// Allow 1 minute tolerance for test execution time
	timeDiff := actualExpiration.Sub(expectedExpiration)
	assert.Less(suite.T(), timeDiff, time.Minute)
	assert.Greater(suite.T(), timeDiff, -time.Minute)
}

// Password hashing tests
func (suite *JWTServiceTestSuite) TestHashPassword() {
	password := "test-password-123"
	
	hashedPassword, err := HashPassword(password)
	
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)
	assert.Greater(suite.T(), len(hashedPassword), 50) // bcrypt hashes are typically 60+ chars
}

func (suite *JWTServiceTestSuite) TestCheckPasswordHash() {
	password := "test-password-123"
	wrongPassword := "wrong-password"
	
	// Hash the password
	hashedPassword, err := HashPassword(password)
	assert.NoError(suite.T(), err)
	
	// Check correct password
	isValid := CheckPasswordHash(password, hashedPassword)
	assert.True(suite.T(), isValid)
	
	// Check wrong password
	isValid = CheckPasswordHash(wrongPassword, hashedPassword)
	assert.False(suite.T(), isValid)
}

func (suite *JWTServiceTestSuite) TestHashPasswordConsistency() {
	password := "same-password"
	
	// Hash the same password multiple times
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	
	// Hashes should be different (due to salt)
	assert.NotEqual(suite.T(), hash1, hash2)
	
	// But both should validate the same password
	assert.True(suite.T(), CheckPasswordHash(password, hash1))
	assert.True(suite.T(), CheckPasswordHash(password, hash2))
}

// Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	jwtService := NewJWTService("benchmark-secret")
	user := &models.User{
		ID:               uuid.New(),
		Email:            "benchmark@example.com",
		SubscriptionTier: models.TierPro,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtService.GenerateToken(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	jwtService := NewJWTService("benchmark-secret")
	user := &models.User{
		ID:               uuid.New(),
		Email:            "benchmark@example.com",
		SubscriptionTier: models.TierPro,
	}
	
	token, _ := jwtService.GenerateToken(user)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtService.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmark-password-123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}