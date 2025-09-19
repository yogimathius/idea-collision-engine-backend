import type { CollisionInput, CollisionResult, CollisionDomain } from './types';
import { collisionDomains } from './collisionDomains';
import { validateCollisionInput, sanitizeUserInput } from './utils/validation';

export class CollisionEngine {
  private generateId(): string {
    return Math.random().toString(36).substr(2, 9);
  }

  async generateCollision(input: CollisionInput): Promise<CollisionResult> {
    // Validate and sanitize input
    const validatedInput = validateCollisionInput({
      ...input,
      currentProject: sanitizeUserInput(input.currentProject),
      userInterests: input.userInterests.map(sanitizeUserInput)
    });

    // 1. Select optimal collision domain
    const selectedDomain = this.selectOptimalDomain(validatedInput);
    
    // 2. Generate connection via OpenAI API
    const connection = await this.generateConnection(validatedInput, selectedDomain);
    
    // 3. Create spark questions
    const sparkQuestions = this.generateSparkQuestions(validatedInput, selectedDomain);
    
    // 4. Generate examples and next steps
    const examples = selectedDomain.examples.slice(0, 3);
    const nextSteps = this.generateNextSteps(validatedInput, selectedDomain);
    
    return {
      id: this.generateId(),
      primaryDomain: this.selectPrimaryDomain(validatedInput),
      collisionDomain: selectedDomain.name,
      connection,
      sparkQuestions,
      examples,
      nextSteps,
      qualityScore: this.calculateQualityScore(validatedInput, selectedDomain),
      timestamp: new Date()
    };
  }

  private selectOptimalDomain(input: CollisionInput): CollisionDomain {
    // Filter domains by intensity
    const validDomains = collisionDomains.filter(domain => 
      domain.intensity.includes(input.collisionIntensity)
    );

    // Score domains based on novelty and relevance
    const scoredDomains = validDomains.map(domain => ({
      domain,
      score: this.scoreDomain(domain, input)
    }));

    // Sort by score and add some randomness
    scoredDomains.sort((a, b) => b.score - a.score);
    
    // Select from top 3 to add variety
    const topDomains = scoredDomains.slice(0, Math.min(3, scoredDomains.length));
    const randomIndex = Math.floor(Math.random() * topDomains.length);
    
    return topDomains[randomIndex].domain;
  }

  private scoreDomain(domain: CollisionDomain, input: CollisionInput): number {
    let score = 0;
    
    // Novelty score - domains less related to user interests score higher
    const interestOverlap = input.userInterests.some(interest => 
      domain.keywords.some(keyword => 
        interest.toLowerCase().includes(keyword.toLowerCase()) ||
        keyword.toLowerCase().includes(interest.toLowerCase())
      )
    );
    
    score += interestOverlap ? 0.3 : 0.8; // Prefer unexpected domains
    
    // Project type relevance
    if (input.projectType === 'product' && ['game_design', 'architecture', 'cooking'].includes(domain.id)) {
      score += 0.3;
    }
    if (input.projectType === 'business' && ['economics', 'ecology', 'game_design'].includes(domain.id)) {
      score += 0.3;
    }
    if (input.projectType === 'research' && ['quantum_physics', 'neuroscience', 'anthropology'].includes(domain.id)) {
      score += 0.3;
    }
    if (input.projectType === 'content' && ['mythology', 'theater', 'music_theory'].includes(domain.id)) {
      score += 0.3;
    }
    
    // Add randomness for variety
    score += Math.random() * 0.2;
    
    return score;
  }

  private selectPrimaryDomain(input: CollisionInput): string {
    return input.userInterests.length > 0 ? input.userInterests[0] : input.projectType;
  }

  private async generateConnection(input: CollisionInput, domain: CollisionDomain): Promise<string> {
    // For MVP, we'll use template-based generation
    // In production, this would call OpenAI API
    const templates = [
      `Just as ${domain.description.toLowerCase()}, your ${input.currentProject} could benefit from similar principles. Consider how ${domain.keywords[0]} and ${domain.keywords[1]} might apply to your ${input.projectType} approach.`,
      
      `The field of ${domain.name.toLowerCase()} reveals that ${domain.keywords[0]} and ${domain.keywords[2]} work together in unexpected ways. This suggests your ${input.currentProject} might explore similar dynamics between seemingly unrelated elements.`,
      
      `${domain.name} shows us that ${domain.keywords[1]} emerges from ${domain.keywords[0]}. What if your ${input.currentProject} approached ${input.userInterests[0] || 'the core challenge'} through this lens of ${domain.keywords[2]}?`
    ];

    return templates[Math.floor(Math.random() * templates.length)];
  }

  private generateSparkQuestions(input: CollisionInput, domain: CollisionDomain): string[] {
    return [
      `What would ${input.currentProject} look like if it followed the principles of ${domain.keywords[0]}?`,
      `How might ${domain.keywords[1]} change your approach to ${input.userInterests[0] || 'this challenge'}?`,
      `Where does your current ${input.projectType} assume linear progress when ${domain.name.toLowerCase()} suggests cyclical or emergent patterns?`,
      `What would happen if you introduced ${domain.keywords[2]} as a core constraint?`,
      `How could the relationship between ${domain.keywords[0]} and ${domain.keywords[1]} inspire a new feature or approach?`
    ].slice(0, 3 + Math.floor(Math.random() * 2));
  }

  private generateNextSteps(input: CollisionInput, domain: CollisionDomain): string[] {
    return [
      `Research how ${domain.keywords[0]} works in ${domain.name.toLowerCase()} and identify 3 principles that could apply`,
      `Experiment with incorporating ${domain.keywords[1]} into your next iteration`,
      `Find experts or communities focused on ${domain.name.toLowerCase()} to learn from`,
      `Create a small prototype that tests ${domain.keywords[2]} in your ${input.projectType} context`
    ].slice(0, 3);
  }

  private calculateQualityScore(input: CollisionInput, domain: CollisionDomain): number {
    // Simple quality scoring for MVP
    let score = 0.7; // Base score
    
    // Bonus for intensity match
    if (domain.intensity.includes(input.collisionIntensity)) {
      score += 0.1;
    }
    
    // Bonus for rich keyword set
    if (domain.keywords.length >= 5) {
      score += 0.1;
    }
    
    // Bonus for good examples
    if (domain.examples.length >= 3) {
      score += 0.1;
    }
    
    return Math.min(1.0, score);
  }
}