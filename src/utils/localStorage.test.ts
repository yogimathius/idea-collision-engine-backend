import { describe, it, expect, beforeEach, vi } from 'vitest';
import { 
  saveCollisionToHistory, 
  getCollisionHistory, 
  updateCollisionRating, 
  addCollisionNotes,
  clearCollisionHistory 
} from './localStorage';
import type { CollisionResult } from '../types';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

// @ts-expect-error - Mocking localStorage for testing
global.localStorage = localStorageMock;

describe('localStorage utilities', () => {
  const mockCollision: CollisionResult = {
    id: 'test-123',
    primaryDomain: 'productivity',
    collisionDomain: 'biomimicry',
    connection: 'Test connection',
    sparkQuestions: ['Test question?'],
    examples: ['Test example'],
    nextSteps: ['Test step'],
    qualityScore: 0.8,
    timestamp: new Date('2024-01-01T00:00:00Z')
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('saveCollisionToHistory', () => {
    it('should save collision to localStorage', () => {
      localStorageMock.getItem.mockReturnValue('[]');
      
      saveCollisionToHistory(mockCollision);
      
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'collision-engine-history',
        expect.stringContaining(mockCollision.id)
      );
    });

    it('should limit history to 50 items', () => {
      const existingHistory = Array.from({ length: 50 }, (_, i) => ({
        ...mockCollision,
        id: `collision-${i}`
      }));
      
      localStorageMock.getItem.mockReturnValue(JSON.stringify(existingHistory));
      
      saveCollisionToHistory(mockCollision);
      
      const savedData = JSON.parse(localStorageMock.setItem.mock.calls[0][1]);
      expect(savedData).toHaveLength(50);
      expect(savedData[0].id).toBe(mockCollision.id);
    });

    it('should handle localStorage errors gracefully', () => {
      localStorageMock.getItem.mockImplementation(() => {
        throw new Error('LocalStorage error');
      });
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      expect(() => saveCollisionToHistory(mockCollision)).not.toThrow();
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });

  describe('getCollisionHistory', () => {
    it('should return empty array when no history exists', () => {
      localStorageMock.getItem.mockReturnValue(null);
      
      const history = getCollisionHistory();
      
      expect(history).toEqual([]);
    });

    it('should parse and return collision history', () => {
      const historyData = [mockCollision];
      localStorageMock.getItem.mockReturnValue(JSON.stringify(historyData));
      
      const history = getCollisionHistory();
      
      expect(history).toHaveLength(1);
      expect(history[0].id).toBe(mockCollision.id);
      expect(history[0].timestamp).toBeInstanceOf(Date);
    });

    it('should handle parse errors gracefully', () => {
      localStorageMock.getItem.mockReturnValue('invalid json');
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      const history = getCollisionHistory();
      
      expect(history).toEqual([]);
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });

  describe('updateCollisionRating', () => {
    it('should update rating for existing collision', () => {
      const historyData = [mockCollision];
      localStorageMock.getItem.mockReturnValue(JSON.stringify(historyData));
      
      updateCollisionRating(mockCollision.id, 5);
      
      const savedData = JSON.parse(localStorageMock.setItem.mock.calls[0][1]);
      expect(savedData[0].rating).toBe(5);
    });

    it('should not affect other collisions', () => {
      const otherCollision = { ...mockCollision, id: 'other-123' };
      const historyData = [mockCollision, otherCollision];
      localStorageMock.getItem.mockReturnValue(JSON.stringify(historyData));
      
      updateCollisionRating(mockCollision.id, 4);
      
      const savedData = JSON.parse(localStorageMock.setItem.mock.calls[0][1]);
      expect(savedData[0].rating).toBe(4);
      expect(savedData[1].rating).toBeUndefined();
    });
  });

  describe('addCollisionNotes', () => {
    it('should add notes to existing collision', () => {
      const historyData = [mockCollision];
      localStorageMock.getItem.mockReturnValue(JSON.stringify(historyData));
      
      const testNotes = 'These are my exploration notes';
      addCollisionNotes(mockCollision.id, testNotes);
      
      const savedData = JSON.parse(localStorageMock.setItem.mock.calls[0][1]);
      expect(savedData[0].notes).toBe(testNotes);
    });
  });

  describe('clearCollisionHistory', () => {
    it('should remove collision history from localStorage', () => {
      clearCollisionHistory();
      
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('collision-engine-history');
    });

    it('should handle clear errors gracefully', () => {
      localStorageMock.removeItem.mockImplementation(() => {
        throw new Error('Clear error');
      });
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      expect(() => clearCollisionHistory()).not.toThrow();
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });
});