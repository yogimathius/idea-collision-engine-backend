import type { CollisionInput } from '../types';

export class ValidationError extends Error {
  public field?: string;
  
  constructor(message: string, field?: string) {
    super(message);
    this.name = 'ValidationError';
    this.field = field;
  }
}

export function validateCollisionInput(input: unknown): CollisionInput {
  if (!input || typeof input !== 'object') {
    throw new ValidationError('Input must be an object');
  }

  const obj = input as Record<string, unknown>;

  // Validate userInterests
  if (!Array.isArray(obj.userInterests)) {
    throw new ValidationError('userInterests must be an array', 'userInterests');
  }

  const userInterests = obj.userInterests.filter(
    (interest): interest is string => 
      typeof interest === 'string' && interest.trim().length > 0
  );

  // Validate currentProject
  if (typeof obj.currentProject !== 'string' || obj.currentProject.trim().length === 0) {
    throw new ValidationError('currentProject is required and must be a non-empty string', 'currentProject');
  }

  const currentProject = obj.currentProject.trim();
  if (currentProject.length < 10) {
    throw new ValidationError('currentProject must be at least 10 characters long', 'currentProject');
  }

  // Validate projectType
  const validProjectTypes = ['product', 'content', 'business', 'research'] as const;
  if (!validProjectTypes.includes(obj.projectType as CollisionInput['projectType'])) {
    throw new ValidationError('projectType must be one of: product, content, business, research', 'projectType');
  }

  // Validate collisionIntensity
  const validIntensities = ['gentle', 'moderate', 'radical'] as const;
  if (!validIntensities.includes(obj.collisionIntensity as CollisionInput['collisionIntensity'])) {
    throw new ValidationError('collisionIntensity must be one of: gentle, moderate, radical', 'collisionIntensity');
  }

  return {
    userInterests,
    currentProject,
    projectType: obj.projectType as CollisionInput['projectType'],
    collisionIntensity: obj.collisionIntensity as CollisionInput['collisionIntensity']
  };
}

export function sanitizeUserInput(input: string): string {
  return input
    .trim()
    .replace(/\s+/g, ' ') // Replace multiple whitespace with single space
    .substring(0, 1000); // Limit length to prevent abuse
}

export function validateRating(rating: unknown): number {
  if (typeof rating !== 'number' || rating < 1 || rating > 5 || !Number.isInteger(rating)) {
    throw new ValidationError('Rating must be an integer between 1 and 5', 'rating');
  }
  return rating;
}

export function validateCollisionId(id: unknown): string {
  if (typeof id !== 'string' || id.length === 0) {
    throw new ValidationError('ID must be a non-empty string', 'id');
  }
  
  // Basic validation for collision ID format
  if (!/^[a-zA-Z0-9_-]+$/.test(id)) {
    throw new ValidationError('ID contains invalid characters', 'id');
  }
  
  return id;
}