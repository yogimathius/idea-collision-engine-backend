import { describe, it, expect } from 'vitest';
import {
  ValidationError,
  validateCollisionInput,
  sanitizeUserInput,
  validateRating,
  validateCollisionId
} from './validation';

describe('validation utilities', () => {
  describe('ValidationError', () => {
    it('should create a ValidationError with message', () => {
      const error = new ValidationError('Test error');
      expect(error.message).toBe('Test error');
      expect(error.name).toBe('ValidationError');
      expect(error instanceof Error).toBe(true);
    });

    it('should create a ValidationError with field', () => {
      const error = new ValidationError('Test error', 'testField');
      expect(error.field).toBe('testField');
    });
  });

  describe('validateCollisionInput', () => {
    const validInput = {
      userInterests: ['productivity', 'design'],
      currentProject: 'Building a task management application for teams',
      projectType: 'product',
      collisionIntensity: 'moderate'
    };

    it('should validate a correct input', () => {
      const result = validateCollisionInput(validInput);
      expect(result).toEqual(validInput);
    });

    it('should throw error for non-object input', () => {
      expect(() => validateCollisionInput(null)).toThrow(ValidationError);
      expect(() => validateCollisionInput('string')).toThrow(ValidationError);
      expect(() => validateCollisionInput(123)).toThrow(ValidationError);
    });

    it('should filter out invalid user interests', () => {
      const input = {
        ...validInput,
        userInterests: ['valid', '', '   ', 'also valid', 123, null]
      };
      
      const result = validateCollisionInput(input);
      expect(result.userInterests).toEqual(['valid', 'also valid']);
    });

    it('should throw error for invalid currentProject', () => {
      expect(() => validateCollisionInput({
        ...validInput,
        currentProject: ''
      })).toThrow(ValidationError);

      expect(() => validateCollisionInput({
        ...validInput,
        currentProject: 'short'
      })).toThrow(ValidationError);

      expect(() => validateCollisionInput({
        ...validInput,
        currentProject: 123
      })).toThrow(ValidationError);
    });

    it('should throw error for invalid projectType', () => {
      expect(() => validateCollisionInput({
        ...validInput,
        projectType: 'invalid'
      })).toThrow(ValidationError);
    });

    it('should throw error for invalid collisionIntensity', () => {
      expect(() => validateCollisionInput({
        ...validInput,
        collisionIntensity: 'extreme'
      })).toThrow(ValidationError);
    });

    it('should trim and sanitize currentProject', () => {
      const input = {
        ...validInput,
        currentProject: '  Building a task management app for teams  '
      };
      
      const result = validateCollisionInput(input);
      expect(result.currentProject).toBe('Building a task management app for teams');
    });
  });

  describe('sanitizeUserInput', () => {
    it('should trim whitespace', () => {
      expect(sanitizeUserInput('  hello world  ')).toBe('hello world');
    });

    it('should replace multiple spaces with single space', () => {
      expect(sanitizeUserInput('hello    world')).toBe('hello world');
      expect(sanitizeUserInput('hello\n\nworld')).toBe('hello world');
      expect(sanitizeUserInput('hello\t\tworld')).toBe('hello world');
    });

    it('should limit length to 1000 characters', () => {
      const longString = 'a'.repeat(1500);
      const result = sanitizeUserInput(longString);
      expect(result.length).toBe(1000);
    });
  });

  describe('validateRating', () => {
    it('should validate correct ratings', () => {
      expect(validateRating(1)).toBe(1);
      expect(validateRating(3)).toBe(3);
      expect(validateRating(5)).toBe(5);
    });

    it('should throw error for invalid ratings', () => {
      expect(() => validateRating(0)).toThrow(ValidationError);
      expect(() => validateRating(6)).toThrow(ValidationError);
      expect(() => validateRating(3.5)).toThrow(ValidationError);
      expect(() => validateRating('3')).toThrow(ValidationError);
      expect(() => validateRating(null)).toThrow(ValidationError);
    });
  });

  describe('validateCollisionId', () => {
    it('should validate correct IDs', () => {
      expect(validateCollisionId('abc123')).toBe('abc123');
      expect(validateCollisionId('test_id-123')).toBe('test_id-123');
    });

    it('should throw error for invalid IDs', () => {
      expect(() => validateCollisionId('')).toThrow(ValidationError);
      expect(() => validateCollisionId('id with spaces')).toThrow(ValidationError);
      expect(() => validateCollisionId('id@with!special')).toThrow(ValidationError);
      expect(() => validateCollisionId(123)).toThrow(ValidationError);
    });
  });
});