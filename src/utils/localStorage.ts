import type { CollisionResult } from '../types';
import { validateCollisionId, validateRating } from './validation';

const COLLISION_HISTORY_KEY = 'collision-engine-history';

export function saveCollisionToHistory(collision: CollisionResult): void {
  try {
    const existingHistory = getCollisionHistory();
    const updatedHistory = [collision, ...existingHistory];
    
    // Keep only the most recent 50 collisions to avoid localStorage bloat
    const trimmedHistory = updatedHistory.slice(0, 50);
    
    localStorage.setItem(COLLISION_HISTORY_KEY, JSON.stringify(trimmedHistory));
  } catch (error) {
    console.error('Failed to save collision to history:', error);
  }
}

export function getCollisionHistory(): CollisionResult[] {
  try {
    const historyJson = localStorage.getItem(COLLISION_HISTORY_KEY);
    if (!historyJson) return [];
    
    const history = JSON.parse(historyJson);
    
    // Convert timestamp strings back to Date objects
    return history.map((collision: CollisionResult & { timestamp: string }) => ({
      ...collision,
      timestamp: new Date(collision.timestamp)
    }));
  } catch (error) {
    console.error('Failed to load collision history:', error);
    return [];
  }
}

export function updateCollisionRating(collisionId: string, rating: number): void {
  try {
    const validatedId = validateCollisionId(collisionId);
    const validatedRating = validateRating(rating);
    
    const history = getCollisionHistory();
    const updatedHistory = history.map(collision =>
      collision.id === validatedId
        ? { ...collision, rating: validatedRating }
        : collision
    );
    
    localStorage.setItem(COLLISION_HISTORY_KEY, JSON.stringify(updatedHistory));
  } catch (error) {
    console.error('Failed to update collision rating:', error);
  }
}

export function addCollisionNotes(collisionId: string, notes: string): void {
  try {
    const validatedId = validateCollisionId(collisionId);
    const sanitizedNotes = notes.trim().substring(0, 2000); // Limit length
    
    const history = getCollisionHistory();
    const updatedHistory = history.map(collision =>
      collision.id === validatedId
        ? { ...collision, notes: sanitizedNotes }
        : collision
    );
    
    localStorage.setItem(COLLISION_HISTORY_KEY, JSON.stringify(updatedHistory));
  } catch (error) {
    console.error('Failed to add collision notes:', error);
  }
}

export function clearCollisionHistory(): void {
  try {
    localStorage.removeItem(COLLISION_HISTORY_KEY);
  } catch (error) {
    console.error('Failed to clear collision history:', error);
  }
}