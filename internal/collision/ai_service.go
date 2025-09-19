package collision

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	"idea-collision-engine-api/internal/models"
)

type AIService struct {
	client *openai.Client
}

func NewAIService(apiKey string) *AIService {
	client := openai.NewClient(apiKey)
	return &AIService{client: client}
}

// EnhanceCollisionResult uses AI to improve the collision with deeper insights
func (ai *AIService) EnhanceCollisionResult(result *models.CollisionResult, input models.CollisionInput, domain models.CollisionDomain) error {
	// Enhance the connection explanation
	enhancedConnection, err := ai.generateEnhancedConnection(result, input, domain)
	if err == nil && enhancedConnection != "" {
		result.Connection = enhancedConnection
	}
	
	// Generate more sophisticated spark questions
	enhancedQuestions, err := ai.generateAdvancedSparkQuestions(input, domain)
	if err == nil && len(enhancedQuestions) > 0 {
		result.SparkQuestions = enhancedQuestions
	}
	
	// Create more contextual examples
	enhancedExamples, err := ai.generateContextualExamples(input, domain)
	if err == nil && len(enhancedExamples) > 0 {
		result.Examples = enhancedExamples
	}
	
	// Generate actionable next steps
	enhancedSteps, err := ai.generateAdvancedNextSteps(input, domain)
	if err == nil && len(enhancedSteps) > 0 {
		result.NextSteps = enhancedSteps
	}
	
	return nil
}

// generateEnhancedConnection creates a deeper explanation of the collision
func (ai *AIService) generateEnhancedConnection(result *models.CollisionResult, input models.CollisionInput, domain models.CollisionDomain) (string, error) {
	prompt := ai.buildConnectionPrompt(input, domain)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 200,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert at finding meaningful connections between disparate fields. Create insightful, practical connections that spark innovation.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}
	
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	
	if len(resp.Choices) > 0 {
		return strings.TrimSpace(resp.Choices[0].Message.Content), nil
	}
	
	return "", fmt.Errorf("no response generated")
}

// generateAdvancedSparkQuestions creates thought-provoking questions
func (ai *AIService) generateAdvancedSparkQuestions(input models.CollisionInput, domain models.CollisionDomain) ([]string, error) {
	prompt := ai.buildSparkQuestionsPrompt(input, domain)
	
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 250,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Generate thought-provoking questions that help people explore unexpected connections. Focus on actionable insights and creative breakthroughs.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.8,
	}
	
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	
	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		questions := ai.parseQuestionsList(content)
		return questions, nil
	}
	
	return nil, fmt.Errorf("no questions generated")
}

// generateContextualExamples creates relevant examples for the specific context
func (ai *AIService) generateContextualExamples(input models.CollisionInput, domain models.CollisionDomain) ([]string, error) {
	prompt := ai.buildExamplesPrompt(input, domain)
	
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 300,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Create specific, actionable examples showing how principles from one domain can be applied to another. Focus on concrete applications.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}
	
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	
	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		examples := ai.parseExamplesList(content)
		return examples, nil
	}
	
	return nil, fmt.Errorf("no examples generated")
}

// generateAdvancedNextSteps creates actionable implementation steps
func (ai *AIService) generateAdvancedNextSteps(input models.CollisionInput, domain models.CollisionDomain) ([]string, error) {
	prompt := ai.buildNextStepsPrompt(input, domain)
	
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 250,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Generate specific, actionable next steps that someone can take to explore and implement cross-domain insights. Be practical and concrete.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.6,
	}
	
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	
	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		steps := ai.parseStepsList(content)
		return steps, nil
	}
	
	return nil, fmt.Errorf("no steps generated")
}

// buildConnectionPrompt constructs the prompt for connection generation
func (ai *AIService) buildConnectionPrompt(input models.CollisionInput, domain models.CollisionDomain) string {
	return fmt.Sprintf(`Create a meaningful connection between %s and "%s" (a %s project).

Domain: %s
Category: %s  
Description: %s
Key concepts: %s

User interests: %s
Collision intensity: %s

Generate a 2-3 sentence explanation of how %s principles can enhance or transform the "%s" project. Focus on specific, actionable insights rather than vague connections.`,
		domain.Name,
		input.CurrentProject,
		input.ProjectType,
		domain.Name,
		domain.Category,
		domain.Description,
		strings.Join(domain.Keywords[:min(5, len(domain.Keywords))], ", "),
		strings.Join(input.UserInterests, ", "),
		input.CollisionIntensity,
		domain.Name,
		input.CurrentProject,
	)
}

// buildSparkQuestionsPrompt constructs the prompt for question generation
func (ai *AIService) buildSparkQuestionsPrompt(input models.CollisionInput, domain models.CollisionDomain) string {
	return fmt.Sprintf(`Generate 4 thought-provoking questions that help someone explore connections between %s and their "%s" project.

Domain: %s
Description: %s
Project type: %s
User interests: %s

Each question should:
- Encourage deep thinking about cross-domain applications
- Be specific and actionable
- Help identify concrete opportunities
- Spark creative breakthroughs

Format as a numbered list (1., 2., 3., 4.).`,
		domain.Name,
		input.CurrentProject,
		domain.Name,
		domain.Description,
		input.ProjectType,
		strings.Join(input.UserInterests, ", "),
	)
}

// buildExamplesPrompt constructs the prompt for example generation
func (ai *AIService) buildExamplesPrompt(input models.CollisionInput, domain models.CollisionDomain) string {
	return fmt.Sprintf(`Generate 3 specific examples showing how %s principles can be applied to a %s project like "%s".

Domain: %s
Description: %s
Key concepts: %s

Each example should:
- Show a specific principle or technique from %s
- Demonstrate concrete application to the %s project
- Be realistic and implementable
- Provide clear value

Format as a numbered list (1., 2., 3.).`,
		domain.Name,
		input.ProjectType,
		input.CurrentProject,
		domain.Name,
		domain.Description,
		strings.Join(domain.Keywords[:min(3, len(domain.Keywords))], ", "),
		domain.Name,
		input.CurrentProject,
	)
}

// buildNextStepsPrompt constructs the prompt for next steps generation
func (ai *AIService) buildNextStepsPrompt(input models.CollisionInput, domain models.CollisionDomain) string {
	return fmt.Sprintf(`Generate 4 actionable next steps for someone wanting to apply %s insights to their "%s" project.

Domain: %s
Project type: %s
User interests: %s

Each step should:
- Be specific and actionable
- Build toward implementing the cross-domain connection
- Be achievable within 1-2 weeks
- Progress from research to implementation

Format as a numbered list (1., 2., 3., 4.).`,
		domain.Name,
		input.CurrentProject,
		domain.Name,
		input.ProjectType,
		strings.Join(input.UserInterests, ", "),
	)
}

// parseQuestionsList extracts questions from AI response
func (ai *AIService) parseQuestionsList(content string) []string {
	return ai.parseNumberedList(content, 4)
}

// parseExamplesList extracts examples from AI response
func (ai *AIService) parseExamplesList(content string) []string {
	return ai.parseNumberedList(content, 3)
}

// parseStepsList extracts steps from AI response
func (ai *AIService) parseStepsList(content string) []string {
	return ai.parseNumberedList(content, 4)
}

// parseNumberedList extracts items from numbered list format
func (ai *AIService) parseNumberedList(content string, expectedCount int) []string {
	lines := strings.Split(content, "\n")
	var items []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Match numbered items (1., 2., etc.)
		for i := 1; i <= expectedCount; i++ {
			prefix := fmt.Sprintf("%d.", i)
			if strings.HasPrefix(line, prefix) {
				item := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				if item != "" {
					items = append(items, item)
				}
				break
			}
		}
	}
	
	return items
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CheckConnection validates OpenAI API connectivity
func (ai *AIService) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 10,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Test connection. Respond with 'OK'.",
			},
		},
	}
	
	_, err := ai.client.CreateChatCompletion(ctx, req)
	return err
}