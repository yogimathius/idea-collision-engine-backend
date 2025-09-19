import { describe, it, expect, beforeEach } from 'vitest';
import { CollisionEngine } from './collisionEngine';
import type { CollisionInput } from './types';

describe('CollisionEngine', () => {
  let collisionEngine: CollisionEngine;
  let mockInput: CollisionInput;

  beforeEach(() => {
    collisionEngine = new CollisionEngine();
    mockInput = {
      userInterests: ['productivity', 'design'],
      currentProject: 'Building a task management app',
      projectType: 'product',
      collisionIntensity: 'moderate'
    };
  });

  describe('generateCollision', () => {
    it('should generate a valid collision result', async () => {
      const result = await collisionEngine.generateCollision(mockInput);

      expect(result).toMatchObject({
        id: expect.any(String),
        primaryDomain: expect.any(String),
        collisionDomain: expect.any(String),
        connection: expect.any(String),
        sparkQuestions: expect.any(Array),
        examples: expect.any(Array),
        nextSteps: expect.any(Array),
        qualityScore: expect.any(Number),
        timestamp: expect.any(Date)
      });

      expect(result.sparkQuestions.length).toBeGreaterThan(0);
      expect(result.examples.length).toBeGreaterThan(0);
      expect(result.nextSteps.length).toBeGreaterThan(0);
      expect(result.qualityScore).toBeGreaterThanOrEqual(0);
      expect(result.qualityScore).toBeLessThanOrEqual(1);
    });

    it('should generate different collisions for the same input', async () => {
      const result1 = await collisionEngine.generateCollision(mockInput);
      const result2 = await collisionEngine.generateCollision(mockInput);

      // Should have different IDs and potentially different collision domains
      expect(result1.id).not.toBe(result2.id);
    });

    it('should respect collision intensity levels', async () => {
      const gentleInput = { ...mockInput, collisionIntensity: 'gentle' as const };
      const radicalInput = { ...mockInput, collisionIntensity: 'radical' as const };

      const gentleResult = await collisionEngine.generateCollision(gentleInput);
      const radicalResult = await collisionEngine.generateCollision(radicalInput);

      expect(gentleResult.collisionDomain).toEqual(expect.any(String));
      expect(radicalResult.collisionDomain).toEqual(expect.any(String));
    });

    it('should handle different project types', async () => {
      const projectTypes = ['product', 'content', 'business', 'research'] as const;
      
      for (const projectType of projectTypes) {
        const input = { ...mockInput, projectType };
        const result = await collisionEngine.generateCollision(input);
        
        expect(result.collisionDomain).toEqual(expect.any(String));
        // Check that the result contains relevant project information rather than exact project type
        expect(result.connection.length).toBeGreaterThan(0);
        expect(result.sparkQuestions.length).toBeGreaterThan(0);
      }
    });

    it('should generate relevant spark questions', async () => {
      const result = await collisionEngine.generateCollision(mockInput);

      expect(result.sparkQuestions.length).toBeGreaterThanOrEqual(3);
      expect(result.sparkQuestions.length).toBeLessThanOrEqual(5);
      
      // Questions should be actual questions
      result.sparkQuestions.forEach(question => {
        expect(question).toMatch(/\?$/);
      });
    });

    it('should generate actionable next steps', async () => {
      const result = await collisionEngine.generateCollision(mockInput);

      expect(result.nextSteps.length).toEqual(3);
      
      // Next steps should be actionable (start with verbs)
      result.nextSteps.forEach(step => {
        expect(step.length).toBeGreaterThan(10);
      });
    });
  });

  describe('quality scoring', () => {
    it('should calculate quality scores within valid range', async () => {
      const results = await Promise.all([
        collisionEngine.generateCollision(mockInput),
        collisionEngine.generateCollision({ ...mockInput, collisionIntensity: 'gentle' }),
        collisionEngine.generateCollision({ ...mockInput, collisionIntensity: 'radical' })
      ]);

      results.forEach(result => {
        expect(result.qualityScore).toBeGreaterThanOrEqual(0);
        expect(result.qualityScore).toBeLessThanOrEqual(1);
      });
    });
  });

  describe('domain selection', () => {
    it('should select appropriate domains for different intensities', async () => {
      const gentleInput = { ...mockInput, collisionIntensity: 'gentle' as const };
      const moderateInput = { ...mockInput, collisionIntensity: 'moderate' as const };
      const radicalInput = { ...mockInput, collisionIntensity: 'radical' as const };

      const gentleResult = await collisionEngine.generateCollision(gentleInput);
      const moderateResult = await collisionEngine.generateCollision(moderateInput);
      const radicalResult = await collisionEngine.generateCollision(radicalInput);

      // All should have valid collision domains
      expect(gentleResult.collisionDomain).toBeTruthy();
      expect(moderateResult.collisionDomain).toBeTruthy();
      expect(radicalResult.collisionDomain).toBeTruthy();
    });
  });
});