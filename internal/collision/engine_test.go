package collision

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"idea-collision-engine-api/internal/models"
)

type CollisionEngineTestSuite struct {
	suite.Suite
	engine *CollisionEngine
	domains []models.CollisionDomain
}

func (suite *CollisionEngineTestSuite) SetupTest() {
	suite.domains = []models.CollisionDomain{
		{
			ID:          uuid.New().String(),
			Name:        "Biomimicry",
			Category:    "Nature",
			Description: "How nature solves problems",
			Keywords:    []string{"evolution", "adaptation", "efficiency"},
			Intensity:   []string{"gentle", "moderate"},
			Tier:        "basic",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Jazz Improvisation",
			Category:    "Music",
			Description: "Spontaneous creation and structured freedom",
			Keywords:    []string{"improvisation", "spontaneity", "collaboration"},
			Intensity:   []string{"moderate", "radical"},
			Tier:        "basic",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Quantum Physics",
			Category:    "Science",
			Description: "Counterintuitive principles of reality",
			Keywords:    []string{"uncertainty", "entanglement", "superposition"},
			Intensity:   []string{"radical"},
			Tier:        "premium",
		},
	}
	
	suite.engine = NewCollisionEngine(suite.domains)
}

func (suite *CollisionEngineTestSuite) TestNewCollisionEngine() {
	engine := NewCollisionEngine(suite.domains)
	
	assert.NotNil(suite.T(), engine)
	assert.Equal(suite.T(), len(suite.domains), len(engine.Domains))
	assert.Equal(suite.T(), suite.domains[0].Name, engine.Domains[0].Name)
}

func (suite *CollisionEngineTestSuite) TestGenerateCollision() {
	input := models.CollisionInput{
		UserInterests:      []string{"machine learning", "design"},
		CurrentProject:     "AI recommendation system",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	result, err := suite.engine.GenerateCollision(input)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.NotEmpty(suite.T(), result.ID)
	assert.NotEmpty(suite.T(), result.PrimaryDomain)
	assert.NotEmpty(suite.T(), result.CollisionDomain)
	assert.NotEmpty(suite.T(), result.Connection)
	assert.GreaterOrEqual(suite.T(), len(result.SparkQuestions), 0)
	assert.GreaterOrEqual(suite.T(), len(result.Examples), 0)
	assert.GreaterOrEqual(suite.T(), len(result.NextSteps), 0)
	assert.GreaterOrEqual(suite.T(), result.QualityScore, 0.0)
	assert.LessOrEqual(suite.T(), result.QualityScore, 100.0)
	assert.WithinDuration(suite.T(), time.Now(), result.Timestamp, 5*time.Second)
}

func (suite *CollisionEngineTestSuite) TestSelectPrimaryDomain() {
	// Test with matching interests
	interests := []string{"nature", "biology"}
	primary := suite.engine.selectPrimaryDomain(interests)
	assert.NotEmpty(suite.T(), primary)
	
	// Test with empty interests
	emptyInterests := []string{}
	primary = suite.engine.selectPrimaryDomain(emptyInterests)
	assert.Equal(suite.T(), "General Innovation", primary)
	
	// Test with non-matching interests
	nonMatchingInterests := []string{"cooking", "photography"}
	primary = suite.engine.selectPrimaryDomain(nonMatchingInterests)
	assert.NotEmpty(suite.T(), primary)
}

func (suite *CollisionEngineTestSuite) TestCalculateInterestRelevance() {
	interests := []string{"nature", "evolution"}
	domain := suite.domains[0] // Biomimicry
	
	relevance := suite.engine.calculateInterestRelevance(interests, domain)
	
	assert.GreaterOrEqual(suite.T(), relevance, 0.0)
	assert.LessOrEqual(suite.T(), relevance, 1.0)
	
	// Test with no interests
	emptyInterests := []string{}
	relevance = suite.engine.calculateInterestRelevance(emptyInterests, domain)
	assert.Equal(suite.T(), 0.0, relevance)
}

func (suite *CollisionEngineTestSuite) TestCalculateNoveltyScore() {
	interests := []string{"machine learning", "technology"}
	domain := suite.domains[1] // Jazz Improvisation - should be novel for tech interests
	
	novelty := suite.engine.calculateNoveltyScore(interests, domain)
	
	assert.GreaterOrEqual(suite.T(), novelty, 0.0)
	assert.LessOrEqual(suite.T(), novelty, 1.0)
	
	// Jazz should be more novel for tech interests than biomimicry
	biomimicryNovelty := suite.engine.calculateNoveltyScore(interests, suite.domains[0])
	assert.GreaterOrEqual(suite.T(), novelty, biomimicryNovelty)
}

func (suite *CollisionEngineTestSuite) TestIsIntensityCompatible() {
	gentle := suite.engine.isIntensityCompatible(suite.domains[0], "gentle") // Biomimicry supports gentle
	assert.True(suite.T(), gentle)
	
	radical := suite.engine.isIntensityCompatible(suite.domains[0], "radical") // Biomimicry doesn't support radical
	assert.False(suite.T(), radical)
	
	quantumRadical := suite.engine.isIntensityCompatible(suite.domains[2], "radical") // Quantum supports radical
	assert.True(suite.T(), quantumRadical)
}

func (suite *CollisionEngineTestSuite) TestFilterCandidateDomains() {
	input := models.CollisionInput{
		UserInterests:      []string{"nature"},
		CurrentProject:     "ecosystem design",
		ProjectType:        "research",
		CollisionIntensity: "gentle",
	}
	
	primaryDomain := "Biomimicry"
	candidates := suite.engine.filterCandidateDomains(input, primaryDomain)
	
	// Should not include the primary domain
	for _, candidate := range candidates {
		assert.NotEqual(suite.T(), primaryDomain, candidate.Name)
	}
	
	// Should only include domains that support the intensity
	for _, candidate := range candidates {
		supported := suite.engine.isIntensityCompatible(candidate, input.CollisionIntensity)
		assert.True(suite.T(), supported, "Domain %s should support intensity %s", candidate.Name, input.CollisionIntensity)
	}
}

func (suite *CollisionEngineTestSuite) TestCalculateQualityScore() {
	input := models.CollisionInput{
		UserInterests:      []string{"artificial intelligence", "innovation"},
		CurrentProject:     "intelligent recommendation system with advanced machine learning algorithms",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	domain := suite.domains[1] // Jazz Improvisation
	score := suite.engine.calculateQualityScore(input, domain)
	
	assert.GreaterOrEqual(suite.T(), score, 0.0)
	assert.LessOrEqual(suite.T(), score, 100.0)
}

func (suite *CollisionEngineTestSuite) TestGenerateConnectionHash() {
	input := models.CollisionInput{
		UserInterests:      []string{"design", "technology"},
		CurrentProject:     "mobile app",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	hash1 := suite.engine.generateConnectionHash(input, "Jazz")
	hash2 := suite.engine.generateConnectionHash(input, "Jazz")
	hash3 := suite.engine.generateConnectionHash(input, "Biomimicry")
	
	// Same input should generate same hash
	assert.Equal(suite.T(), hash1, hash2)
	
	// Different domain should generate different hash
	assert.NotEqual(suite.T(), hash1, hash3)
	
	// Hash should be non-empty and reasonable length
	assert.NotEmpty(suite.T(), hash1)
	assert.Equal(suite.T(), 16, len(hash1)) // First 16 chars of SHA256
}

func (suite *CollisionEngineTestSuite) TestDifferentIntensityLevels() {
	input := models.CollisionInput{
		UserInterests:      []string{"business", "innovation"},
		CurrentProject:     "startup idea",
		ProjectType:        "business",
		CollisionIntensity: "gentle",
	}
	
	// Test gentle intensity
	result, err := suite.engine.GenerateCollision(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	
	// Test moderate intensity
	input.CollisionIntensity = "moderate"
	result, err = suite.engine.GenerateCollision(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	
	// Test radical intensity (should only get quantum physics)
	input.CollisionIntensity = "radical"
	result, err = suite.engine.GenerateCollision(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *CollisionEngineTestSuite) TestAntiEchoChamberScoring() {
	relevantInterests := []string{"biology", "nature", "evolution"}
	
	biomimicryDomain := suite.domains[0] // Should be highly relevant to biology interests
	jazzDomain := suite.domains[1]       // Should be novel for biology interests
	
	// Test relevance calculation
	highRelevance := suite.engine.calculateInterestRelevance(relevantInterests, biomimicryDomain)
	lowRelevance := suite.engine.calculateInterestRelevance(relevantInterests, jazzDomain)
	assert.Greater(suite.T(), highRelevance, lowRelevance)
	
	// Test novelty calculation (inverse relationship)
	lowNovelty := suite.engine.calculateNoveltyScore(relevantInterests, biomimicryDomain)
	highNovelty := suite.engine.calculateNoveltyScore(relevantInterests, jazzDomain)
	assert.Greater(suite.T(), highNovelty, lowNovelty)
	
	// Radical should favor novelty more than gentle
	radicalNovelScore := suite.engine.calculateAntiEchoChamberScore(lowRelevance, highNovelty, "radical")
	gentleNovelScore := suite.engine.calculateAntiEchoChamberScore(lowRelevance, highNovelty, "gentle")
	
	assert.Greater(suite.T(), radicalNovelScore, gentleNovelScore)
}

// Benchmark tests
func BenchmarkGenerateCollision(b *testing.B) {
	domains := []models.CollisionDomain{
		{
			ID:          uuid.New().String(),
			Name:        "Test Domain",
			Category:    "Test",
			Description: "Test description",
			Keywords:    []string{"test", "benchmark"},
			Intensity:   []string{"moderate"},
			Tier:        "basic",
		},
	}
	
	engine := NewCollisionEngine(domains)
	input := models.CollisionInput{
		UserInterests:      []string{"technology", "innovation"},
		CurrentProject:     "test project",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.GenerateCollision(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test edge cases
func (suite *CollisionEngineTestSuite) TestEdgeCases() {
	// Test with empty domains
	emptyEngine := NewCollisionEngine([]models.CollisionDomain{})
	input := models.CollisionInput{
		UserInterests:      []string{"test"},
		CurrentProject:     "test project",
		ProjectType:        "product",
		CollisionIntensity: "moderate",
	}
	
	result, err := emptyEngine.GenerateCollision(input)
	
	// Should handle gracefully with fallback
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Innovation", result.CollisionDomain)
}

func TestCollisionEngineTestSuite(t *testing.T) {
	suite.Run(t, new(CollisionEngineTestSuite))
}