package collision

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"idea-collision-engine-api/internal/models"
)

type CollisionEngine struct {
	Domains []models.CollisionDomain
}

type DomainMatch struct {
	Domain        models.CollisionDomain
	RelevanceScore float64
	NoveltyScore   float64
	OverallScore   float64
	Reasoning      string
}

func NewCollisionEngine(domains []models.CollisionDomain) *CollisionEngine {
	return &CollisionEngine{
		Domains: domains,
	}
}

// GenerateCollision creates a collision between user interests and an unexpected domain
func (e *CollisionEngine) GenerateCollision(input models.CollisionInput) (*models.CollisionResult, error) {
	// 1. Find primary domain from user interests
	primaryDomain := e.selectPrimaryDomain(input.UserInterests)
	
	// 2. Apply anti-echo chamber algorithm to find collision domain
	collisionDomain, reasoning := e.selectCollisionDomain(input, primaryDomain)
	
	// 3. Generate connection hash for caching (for future use)
	_ = e.generateConnectionHash(input, collisionDomain.Name)
	
	// 4. Create collision result structure
	result := &models.CollisionResult{
		ID:              uuid.New().String(),
		PrimaryDomain:   primaryDomain,
		CollisionDomain: collisionDomain.Name,
		Connection:      reasoning,
		QualityScore:    e.calculateQualityScore(input, collisionDomain),
		Timestamp:       time.Now(),
	}
	
	// 5. Generate spark questions, examples, and next steps
	e.enrichCollisionResult(result, input, collisionDomain)
	
	return result, nil
}

// selectPrimaryDomain chooses the most relevant domain from user interests
func (e *CollisionEngine) selectPrimaryDomain(interests []string) string {
	if len(interests) == 0 {
		return "General Innovation"
	}
	
	// Find the domain that best matches user interests
	bestMatch := ""
	highestScore := 0.0
	
	for _, domain := range e.Domains {
		score := e.calculateInterestRelevance(interests, domain)
		if score > highestScore {
			highestScore = score
			bestMatch = domain.Name
		}
	}
	
	if bestMatch == "" {
		// Fallback to first interest if no domain matches
		return strings.Title(interests[0])
	}
	
	return bestMatch
}

// selectCollisionDomain implements anti-echo chamber algorithm
func (e *CollisionEngine) selectCollisionDomain(input models.CollisionInput, primaryDomain string) (models.CollisionDomain, string) {
	candidates := e.filterCandidateDomains(input, primaryDomain)
	
	// Score each candidate domain
	var matches []DomainMatch
	for _, domain := range candidates {
		relevance := e.calculateDomainRelevance(input, domain)
		novelty := e.calculateNoveltyScore(input.UserInterests, domain)
		
		// Anti-echo chamber: prioritize novelty while maintaining some relevance
		overall := e.calculateAntiEchoChamberScore(relevance, novelty, input.CollisionIntensity)
		
		reasoning := e.generateReasoningSnippet(input.CurrentProject, domain, relevance, novelty)
		
		matches = append(matches, DomainMatch{
			Domain:        domain,
			RelevanceScore: relevance,
			NoveltyScore:   novelty,
			OverallScore:   overall,
			Reasoning:      reasoning,
		})
	}
	
	// Sort by overall score and add randomness
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].OverallScore > matches[j].OverallScore
	})
	
	// Select from top candidates with weighted randomness
	selected := e.selectWithRandomness(matches, input.CollisionIntensity)
	return selected.Domain, selected.Reasoning
}

// filterCandidateDomains removes unsuitable domains
func (e *CollisionEngine) filterCandidateDomains(input models.CollisionInput, primaryDomain string) []models.CollisionDomain {
	var candidates []models.CollisionDomain
	
	for _, domain := range e.Domains {
		// Skip if it's the same as primary domain
		if domain.Name == primaryDomain {
			continue
		}
		
		// Check intensity compatibility
		if !e.isIntensityCompatible(domain, input.CollisionIntensity) {
			continue
		}
		
		candidates = append(candidates, domain)
	}
	
	return candidates
}

// calculateInterestRelevance scores how well a domain matches user interests
func (e *CollisionEngine) calculateInterestRelevance(interests []string, domain models.CollisionDomain) float64 {
	if len(interests) == 0 {
		return 0.0
	}
	
	score := 0.0
	totalPossible := 0.0
	
	for _, interest := range interests {
		interestLower := strings.ToLower(interest)
		domainScore := 0.0
		
		// Check name match
		if strings.Contains(strings.ToLower(domain.Name), interestLower) {
			domainScore += 3.0
		}
		
		// Check category match
		if strings.Contains(strings.ToLower(domain.Category), interestLower) {
			domainScore += 2.0
		}
		
		// Check keyword matches
		for _, keyword := range domain.Keywords {
			if strings.Contains(strings.ToLower(keyword), interestLower) ||
				strings.Contains(interestLower, strings.ToLower(keyword)) {
				domainScore += 1.0
			}
		}
		
		// Check description match
		if strings.Contains(strings.ToLower(domain.Description), interestLower) {
			domainScore += 0.5
		}
		
		score += math.Min(domainScore, 3.0) // Cap individual interest score
		totalPossible += 3.0
	}
	
	return score / totalPossible
}

// calculateDomainRelevance scores domain relevance to project context
func (e *CollisionEngine) calculateDomainRelevance(input models.CollisionInput, domain models.CollisionDomain) float64 {
	score := 0.0
	
	projectLower := strings.ToLower(input.CurrentProject)
	projectTypeLower := strings.ToLower(input.ProjectType)
	
	// Project type relevance
	categoryLower := strings.ToLower(domain.Category)
	
	relevanceMap := map[string][]string{
		"product":  {"design", "technology", "science", "crafts"},
		"content":  {"arts", "media", "cultural", "entertainment"},
		"business": {"social systems", "economics", "human systems"},
		"research": {"science", "mathematics", "philosophy"},
	}
	
	if categories, exists := relevanceMap[projectTypeLower]; exists {
		for _, cat := range categories {
			if strings.Contains(categoryLower, cat) {
				score += 0.3
				break
			}
		}
	}
	
	// Project description relevance
	for _, keyword := range domain.Keywords {
		if strings.Contains(projectLower, strings.ToLower(keyword)) {
			score += 0.2
		}
	}
	
	// Example relevance
	for _, example := range domain.Examples {
		exampleLower := strings.ToLower(example)
		words := strings.Fields(projectLower)
		for _, word := range words {
			if len(word) > 3 && strings.Contains(exampleLower, word) {
				score += 0.1
				break
			}
		}
	}
	
	return math.Min(score, 1.0)
}

// calculateNoveltyScore measures how unexpected the domain is
func (e *CollisionEngine) calculateNoveltyScore(interests []string, domain models.CollisionDomain) float64 {
	// Higher novelty = lower relevance to existing interests
	relevance := e.calculateInterestRelevance(interests, domain)
	
	// Invert relevance for novelty, but keep some floor
	novelty := math.Max(0.2, 1.0-relevance)
	
	// Boost novelty for certain categories that are inherently unexpected
	unexpectedCategories := []string{"quantum", "chaos", "mythology", "ancient", "radical"}
	categoryLower := strings.ToLower(domain.Category + " " + domain.Name + " " + domain.Description)
	
	for _, unexpected := range unexpectedCategories {
		if strings.Contains(categoryLower, unexpected) {
			novelty *= 1.2
			break
		}
	}
	
	return math.Min(novelty, 1.0)
}

// calculateAntiEchoChamberScore balances relevance and novelty
func (e *CollisionEngine) calculateAntiEchoChamberScore(relevance, novelty float64, intensity string) float64 {
	// Weight novelty higher to break echo chambers
	weights := map[string][2]float64{
		"gentle":   {0.6, 0.4}, // 60% relevance, 40% novelty
		"moderate": {0.4, 0.6}, // 40% relevance, 60% novelty
		"radical":  {0.2, 0.8}, // 20% relevance, 80% novelty
	}
	
	weight, exists := weights[intensity]
	if !exists {
		weight = weights["moderate"]
	}
	
	return relevance*weight[0] + novelty*weight[1]
}

// selectWithRandomness adds controlled randomness to selection
func (e *CollisionEngine) selectWithRandomness(matches []DomainMatch, intensity string) DomainMatch {
	if len(matches) == 0 {
		// Fallback - this shouldn't happen
		return DomainMatch{
			Domain: models.CollisionDomain{
				Name:        "Innovation",
				Category:    "General",
				Description: "General innovative thinking",
			},
			Reasoning: "Fallback domain for creative exploration",
		}
	}
	
	// Define selection pool size based on intensity
	poolSizes := map[string]int{
		"gentle":   3, // Choose from top 3
		"moderate": 5, // Choose from top 5
		"radical":  8, // Choose from top 8
	}
	
	poolSize := poolSizes["moderate"]
	if size, exists := poolSizes[intensity]; exists {
		poolSize = size
	}
	
	if poolSize > len(matches) {
		poolSize = len(matches)
	}
	
	// Use weighted randomness - higher scores more likely
	weights := make([]float64, poolSize)
	totalWeight := 0.0
	
	for i := 0; i < poolSize; i++ {
		// Exponential decay for weighting
		weights[i] = math.Exp(-float64(i) * 0.5)
		totalWeight += weights[i]
	}
	
	// Select randomly based on weights
	rand.Seed(time.Now().UnixNano())
	target := rand.Float64() * totalWeight
	
	cumulative := 0.0
	for i := 0; i < poolSize; i++ {
		cumulative += weights[i]
		if cumulative >= target {
			return matches[i]
		}
	}
	
	// Fallback to first match
	return matches[0]
}

// isIntensityCompatible checks if domain supports the requested intensity
func (e *CollisionEngine) isIntensityCompatible(domain models.CollisionDomain, intensity string) bool {
	for _, supportedIntensity := range domain.Intensity {
		if supportedIntensity == intensity {
			return true
		}
	}
	return false
}

// calculateQualityScore provides overall collision quality assessment
func (e *CollisionEngine) calculateQualityScore(input models.CollisionInput, domain models.CollisionDomain) float64 {
	relevance := e.calculateDomainRelevance(input, domain)
	novelty := e.calculateNoveltyScore(input.UserInterests, domain)
	
	// Additional factors
	projectComplexity := e.assessProjectComplexity(input.CurrentProject)
	domainDepth := e.assessDomainDepth(domain)
	
	// Weighted average
	score := (relevance*0.3 + novelty*0.3 + projectComplexity*0.2 + domainDepth*0.2) * 100
	
	// Add some randomness to prevent identical scores
	rand.Seed(time.Now().UnixNano())
	score += (rand.Float64() - 0.5) * 5 // ±2.5 points
	
	return math.Max(0, math.Min(100, score))
}

// assessProjectComplexity estimates project sophistication
func (e *CollisionEngine) assessProjectComplexity(project string) float64 {
	complexityIndicators := []string{
		"system", "platform", "algorithm", "network", "framework",
		"architecture", "optimization", "intelligence", "automation",
		"integration", "scalable", "distributed", "analytics",
	}
	
	projectLower := strings.ToLower(project)
	matches := 0
	
	for _, indicator := range complexityIndicators {
		if strings.Contains(projectLower, indicator) {
			matches++
		}
	}
	
	// Normalize to 0-1 range
	return math.Min(1.0, float64(matches)/5.0)
}

// assessDomainDepth evaluates domain sophistication
func (e *CollisionEngine) assessDomainDepth(domain models.CollisionDomain) float64 {
	score := 0.0
	
	// More keywords = higher depth
	score += math.Min(0.3, float64(len(domain.Keywords))/10.0)
	
	// More examples = higher depth
	score += math.Min(0.3, float64(len(domain.Examples))/5.0)
	
	// Description length indicates depth
	score += math.Min(0.2, float64(len(domain.Description))/200.0)
	
	// Premium tier = higher depth
	if domain.Tier == "premium" {
		score += 0.2
	}
	
	return math.Min(1.0, score)
}

// generateReasoningSnippet creates the initial connection explanation
func (e *CollisionEngine) generateReasoningSnippet(project string, domain models.CollisionDomain, relevance, novelty float64) string {
	if novelty > 0.7 {
		return fmt.Sprintf("Exploring %s offers an unexpected lens for %s, challenging conventional approaches through %s principles.",
			domain.Name, project, strings.ToLower(domain.Category))
	} else if relevance > 0.6 {
		return fmt.Sprintf("The principles of %s can directly enhance %s by applying %s methodologies.",
			domain.Name, project, strings.ToLower(domain.Category))
	} else {
		return fmt.Sprintf("Drawing from %s creates novel opportunities for %s through cross-disciplinary insight.",
			domain.Name, project)
	}
}

// enrichCollisionResult adds spark questions, examples, and next steps
func (e *CollisionEngine) enrichCollisionResult(result *models.CollisionResult, input models.CollisionInput, domain models.CollisionDomain) {
	// Generate spark questions
	result.SparkQuestions = e.generateSparkQuestions(input, domain)
	
	// Adapt domain examples to the specific project
	result.Examples = e.adaptExamples(input, domain)
	
	// Create actionable next steps
	result.NextSteps = e.generateNextSteps(input, domain)
}

// generateSparkQuestions creates thought-provoking questions
func (e *CollisionEngine) generateSparkQuestions(input models.CollisionInput, domain models.CollisionDomain) []string {
	questions := []string{
		fmt.Sprintf("How might %s principles reshape your approach to %s?",
			strings.ToLower(domain.Name), input.CurrentProject),
		fmt.Sprintf("What would %s look like if designed using %s patterns?",
			input.CurrentProject, domain.Category),
		fmt.Sprintf("Which aspects of %s could introduce unexpected benefits to your %s project?",
			domain.Name, input.ProjectType),
	}
	
	// Add domain-specific questions based on keywords
	if len(domain.Keywords) > 0 {
		keyword := domain.Keywords[rand.Intn(len(domain.Keywords))]
		questions = append(questions,
			fmt.Sprintf("How could the concept of '%s' unlock new possibilities in your work?", keyword))
	}
	
	return questions
}

// adaptExamples contextualizes domain examples for the project
func (e *CollisionEngine) adaptExamples(input models.CollisionInput, domain models.CollisionDomain) []string {
	adapted := make([]string, 0, len(domain.Examples))
	
	for _, example := range domain.Examples {
		// Try to contextualize each example
		contextualizedExample := fmt.Sprintf("%s → Applied to %s: %s",
			example, input.ProjectType, e.contextualizeExample(example, input))
		adapted = append(adapted, contextualizedExample)
	}
	
	return adapted
}

// contextualizeExample adapts a domain example to the specific project context
func (e *CollisionEngine) contextualizeExample(example string, input models.CollisionInput) string {
	exampleLower := strings.ToLower(example)
	
	// Simple pattern matching for contextualization
	if strings.Contains(exampleLower, "system") {
		return "could inspire new system architectures"
	} else if strings.Contains(exampleLower, "pattern") {
		return "might reveal new design patterns"
	} else if strings.Contains(exampleLower, "flow") {
		return "could optimize process flows"
	} else {
		return "offers fresh perspective on implementation"
	}
}

// generateNextSteps creates actionable recommendations
func (e *CollisionEngine) generateNextSteps(input models.CollisionInput, domain models.CollisionDomain) []string {
	steps := []string{
		fmt.Sprintf("Research core %s principles and identify 3 that could apply to %s",
			domain.Name, input.CurrentProject),
		fmt.Sprintf("Find experts or resources in %s to deepen understanding",
			domain.Name),
		fmt.Sprintf("Prototype one small aspect of %s using %s-inspired approaches",
			input.CurrentProject, domain.Name),
		fmt.Sprintf("Document insights and unexpected connections discovered"),
	}
	
	// Add intensity-specific steps
	if input.CollisionIntensity == "radical" {
		steps = append(steps,
			fmt.Sprintf("Challenge fundamental assumptions about %s using %s perspective",
				input.ProjectType, domain.Name))
	}
	
	return steps
}

// generateConnectionHash creates a hash for caching similar collision requests
func (e *CollisionEngine) generateConnectionHash(input models.CollisionInput, domainName string) string {
	content := fmt.Sprintf("%v|%s|%s|%s|%s",
		input.UserInterests,
		input.CurrentProject,
		input.ProjectType,
		input.CollisionIntensity,
		domainName)
	
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)[:16] // First 16 chars
}