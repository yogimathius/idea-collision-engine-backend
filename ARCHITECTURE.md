# Idea Collision Engine - Technical Architecture

## Technology Stack Decisions

### Frontend Stack
**React + TypeScript + Tailwind CSS**

**Decision Rationale**:
- **React**: Proven component-based architecture, excellent for interactive UIs
- **TypeScript**: Type safety critical for complex collision data structures
- **Tailwind CSS**: Rapid UI development, consistent design system
- **Vite**: Fast development server, modern build tooling
- **TanStack Query**: Optimistic updates, caching, background sync for collision data

**Key Libraries**:
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "typescript": "^5.0.0",
    "tailwindcss": "^3.3.0",
    "@tanstack/react-query": "^5.0.0",
    "react-router-dom": "^6.8.0",
    "framer-motion": "^10.0.0",
    "zustand": "^4.4.0"
  }
}
```

### Backend Stack
**Node.js + Express + TypeScript**

**Decision Rationale**:
- **Node.js**: JavaScript ecosystem consistency, excellent AI API integration
- **Express**: Lightweight, flexible, extensive middleware ecosystem
- **TypeScript**: Shared types with frontend, better maintainability
- **Prisma**: Type-safe database access, excellent PostgreSQL support
- **Socket.io**: Real-time collaboration for team sessions

**Key Libraries**:
```json
{
  "dependencies": {
    "express": "^4.18.0",
    "typescript": "^5.0.0",
    "@prisma/client": "^5.0.0",
    "socket.io": "^4.7.0",
    "jsonwebtoken": "^9.0.0",
    "openai": "^4.0.0",
    "helmet": "^7.0.0",
    "cors": "^2.8.5"
  }
}
```

### Database Architecture
**PostgreSQL + Prisma ORM**

**Decision Rationale**:
- **PostgreSQL**: JSONB support for flexible collision data, excellent performance
- **Prisma**: Type-safe queries, excellent migration system, great TypeScript integration
- **Redis**: Session storage and caching for performance

## System Architecture

### High-Level Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Client  │    │  Express API    │    │  PostgreSQL DB  │
│                 │    │                 │    │                 │
│ • Collision UI  │◄──►│ • Collision     │◄──►│ • Users         │
│ • History Views │    │   Generation    │    │ • Sessions      │
│ • Team Sessions │    │ • Pattern       │    │ • Domains       │
│                 │    │   Analysis      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐
                       │   OpenAI API    │
                       │                 │
                       │ • Connection    │
                       │   Generation    │
                       │ • Enhancement   │
                       └─────────────────┘
```

### Core Data Models

```typescript
// Core collision types
interface CollisionInput {
  userInterests: string[];
  currentProject: string;
  projectType: 'product' | 'content' | 'business' | 'research';
  collisionIntensity: 'gentle' | 'moderate' | 'radical';
}

interface CollisionResult {
  id: string;
  primaryDomain: string;
  collisionDomain: string;
  connection: string;
  sparkQuestions: string[];
  examples: string[];
  nextSteps: string[];
  qualityScore: number;
}

// Database schema types
interface User {
  id: string;
  email: string;
  subscriptionTier: 'free' | 'pro' | 'team' | 'enterprise';
  interests: string[];
  createdAt: Date;
}

interface CollisionSession {
  id: string;
  userId: string;
  input: CollisionInput;
  result: CollisionResult;
  userRating: number | null;
  explorationNotes: string | null;
  createdAt: Date;
}

interface CollisionDomain {
  id: string;
  name: string;
  category: string;
  description: string;
  exampleApplications: string[];
  tier: 'basic' | 'pro' | 'custom';
  usageCount: number;
}
```

### Collision Generation Algorithm

```typescript
class CollisionEngine {
  async generateCollision(input: CollisionInput): Promise<CollisionResult> {
    // 1. Domain Selection Strategy
    const availableDomains = await this.getAvailableDomains(
      input.projectType,
      input.collisionIntensity
    );
    
    const selectedDomain = await this.selectOptimalDomain(
      input.userInterests,
      availableDomains,
      input.collisionIntensity
    );
    
    // 2. Connection Generation
    const connection = await this.generateConnection(
      input,
      selectedDomain
    );
    
    // 3. Enhancement & Enrichment
    const enhanced = await this.enrichCollision(
      connection,
      selectedDomain,
      input
    );
    
    // 4. Quality Scoring
    const qualityScore = this.calculateQualityScore(enhanced, input);
    
    return enhanced;
  }
  
  private async selectOptimalDomain(
    userInterests: string[],
    domains: CollisionDomain[],
    intensity: string
  ): Promise<CollisionDomain> {
    // Weighted selection based on:
    // - Novelty (how unexpected for this user)
    // - Relevance (connection potential)
    // - Intensity matching
    // - Usage frequency (avoid overused domains)
  }
  
  private calculateQualityScore(
    result: CollisionResult,
    input: CollisionInput
  ): number {
    return (
      this.scoreRelevance(result, input) * 0.3 +
      this.scoreNovelty(result, input) * 0.3 +
      this.scoreActionability(result) * 0.2 +
      this.scoreDepth(result) * 0.2
    );
  }
}
```

## API Design

### REST Endpoints
```typescript
// Authentication
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh
DELETE /api/auth/logout

// Collision Generation
POST   /api/collisions/generate
GET    /api/collisions/history
PUT    /api/collisions/:id/rating
PUT    /api/collisions/:id/notes

// Pattern Analysis
GET    /api/patterns/user-preferences
GET    /api/patterns/collision-evolution/:sessionId
GET    /api/patterns/serendipity-suggestions

// Team Collaboration
POST   /api/teams/sessions
GET    /api/teams/sessions/:id
PUT    /api/teams/sessions/:id/participant
PUT    /api/teams/sessions/:id/synthesis

// Subscription & Usage
GET    /api/user/subscription
GET    /api/user/usage
POST   /api/subscription/upgrade
```

### WebSocket Events (Team Sessions)
```typescript
// Real-time collaboration events
interface TeamSessionEvents {
  'session:join': { sessionId: string; userId: string };
  'session:collision-share': { collision: CollisionResult };
  'session:synthesis-update': { synthesis: SessionSynthesis };
  'session:participant-update': { participants: TeamMember[] };
}
```

## Database Schema

```sql
-- Core user management
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  subscription_tier VARCHAR(20) DEFAULT 'free',
  interests JSONB DEFAULT '[]',
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Collision sessions and results
CREATE TABLE collision_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  input_data JSONB NOT NULL,
  collision_result JSONB NOT NULL,
  user_rating INTEGER CHECK (user_rating >= 1 AND user_rating <= 5),
  exploration_notes TEXT,
  quality_score DECIMAL(3,2),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Collision domain database
CREATE TABLE collision_domains (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  category VARCHAR(50) NOT NULL,
  description TEXT NOT NULL,
  example_applications JSONB DEFAULT '[]',
  tier VARCHAR(20) DEFAULT 'basic',
  usage_count INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Team collaboration sessions
CREATE TABLE team_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  facilitator_id UUID NOT NULL REFERENCES users(id),
  session_type VARCHAR(50) NOT NULL,
  shared_context JSONB NOT NULL,
  synthesis_data JSONB,
  status VARCHAR(20) DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE team_session_participants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES team_sessions(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id),
  joined_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(session_id, user_id)
);

-- Usage tracking for freemium limits
CREATE TABLE user_usage (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  week_start DATE NOT NULL,
  collision_count INTEGER DEFAULT 0,
  team_session_count INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, week_start)
);

-- Indexes for performance
CREATE INDEX idx_collision_sessions_user_created ON collision_sessions(user_id, created_at DESC);
CREATE INDEX idx_collision_sessions_rating ON collision_sessions(user_rating) WHERE user_rating IS NOT NULL;
CREATE INDEX idx_domains_tier_category ON collision_domains(tier, category);
CREATE INDEX idx_user_usage_week ON user_usage(user_id, week_start);
```

## Performance Considerations

### Caching Strategy
```typescript
// Redis caching layers
interface CacheStrategy {
  // Domain data (rarely changes)
  collisionDomains: '24h';
  
  // User patterns (computed expensive)
  userPatterns: '1h';
  
  // API responses (OpenAI calls)
  connectionGeneration: '15min';
  
  // Session data
  userSessions: '30min';
}
```

### Optimization Approaches
1. **Database Query Optimization**
   - Proper indexing on frequently queried fields
   - JSONB queries optimization for collision data
   - Connection pooling for high concurrency

2. **AI API Optimization**
   - Response caching for similar collision requests
   - Batch processing for team sessions
   - Fallback to template-based generation if API fails

3. **Frontend Performance**
   - Code splitting by routes
   - Lazy loading of collision history
   - Optimistic updates for better UX

## Security Architecture

### Authentication & Authorization
```typescript
// JWT token strategy
interface TokenPayload {
  userId: string;
  subscriptionTier: string;
  iat: number;
  exp: number;
}

// Rate limiting by subscription tier
const rateLimits = {
  free: { collisionsPerHour: 5 },
  pro: { collisionsPerHour: 50 },
  team: { collisionsPerHour: 200 },
  enterprise: { collisionsPerHour: -1 }
};
```

### Data Protection
- Helmet.js for security headers
- CORS configuration for frontend domain
- Input validation with Joi schemas
- SQL injection prevention with Prisma
- Environment variable management for secrets

## Deployment Architecture

### Development Environment
```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: collision_engine
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
  
  redis:
    image: redis:7
    ports:
      - "6379:6379"
  
  app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - postgres
      - redis
```

### Production Deployment
- **Frontend**: Vercel (React deployment, CDN, automatic deployments)
- **Backend**: Railway (Node.js + PostgreSQL, automatic scaling)
- **Redis**: Railway Redis add-on for session storage
- **Monitoring**: Sentry for error tracking, PostHog for analytics

## Scalability Planning

### Phase 1: Single Server (0-1K users)
- Single Express server
- Single PostgreSQL instance
- Basic Redis caching

### Phase 2: Horizontal Scaling (1K-10K users)
- Load balancer (Railway automatic)
- Read replicas for database
- Enhanced caching strategy

### Phase 3: Microservices (10K+ users)
- Collision service separation
- Pattern analysis service
- Team collaboration service
- Event-driven architecture with message queues

## Monitoring & Analytics

### Key Metrics to Track
```typescript
interface AnalyticsEvents {
  'collision_generated': {
    userId: string;
    qualityScore: number;
    collisionIntensity: string;
    processingTime: number;
  };
  
  'collision_rated': {
    userId: string;
    sessionId: string;
    rating: number;
    previousRatings: number[];
  };
  
  'pattern_insight_viewed': {
    userId: string;
    patternType: string;
    actionTaken: boolean;
  };
  
  'subscription_upgraded': {
    userId: string;
    fromTier: string;
    toTier: string;
    triggerFeature: string;
  };
}
```

### Health Checks
- API response times
- Database query performance
- OpenAI API success rates
- User satisfaction scores (collision ratings)
- Conversion funnel metrics

This architecture provides a solid foundation for the Idea Collision Engine while maintaining flexibility for future enhancements and scaling needs.