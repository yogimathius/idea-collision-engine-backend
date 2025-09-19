export interface CollisionInput {
  userInterests: string[];
  currentProject: string;
  projectType: 'product' | 'content' | 'business' | 'research';
  collisionIntensity: 'gentle' | 'moderate' | 'radical';
}

export interface CollisionResult {
  id: string;
  primaryDomain: string;
  collisionDomain: string;
  connection: string;
  sparkQuestions: string[];
  examples: string[];
  nextSteps: string[];
  qualityScore: number;
  timestamp: Date;
  rating?: number;
  notes?: string;
}

export interface CollisionDomain {
  id: string;
  name: string;
  category: string;
  description: string;
  examples: string[];
  keywords: string[];
  intensity: ('gentle' | 'moderate' | 'radical')[];
}