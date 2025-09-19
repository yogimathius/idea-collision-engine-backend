import type { CollisionDomain } from './types';

export const collisionDomains: CollisionDomain[] = [
  {
    id: 'biomimicry',
    name: 'Biomimicry',
    category: 'Nature & Biology',
    description: 'How nature solves similar problems through millions of years of evolution',
    examples: ['Velcro from burdock burrs', 'Bullet train design from kingfisher beaks', 'Self-healing materials from plant tissues'],
    keywords: ['evolution', 'adaptation', 'efficiency', 'sustainability', 'natural selection'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'ancient_civilizations',
    name: 'Ancient Civilizations',
    category: 'Historical',
    description: 'Time-tested approaches and wisdom from past cultures',
    examples: ['Roman concrete lasting 2000 years', 'Incan agricultural terraces', 'Greek philosophical methods'],
    keywords: ['wisdom', 'durability', 'systems', 'culture', 'timeless'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'game_design',
    name: 'Game Design',
    category: 'Entertainment',
    description: 'Engagement mechanics, progression systems, and motivation psychology',
    examples: ['Leveling systems for skill building', 'Achievement badges for motivation', 'Balancing challenge and reward'],
    keywords: ['engagement', 'progression', 'feedback', 'motivation', 'flow state'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'music_theory',
    name: 'Music Theory',
    category: 'Arts',
    description: 'Harmony, dissonance, rhythm, and emotional resonance principles',
    examples: ['Tension and resolution in storytelling', 'Rhythmic patterns in UI design', 'Harmonic collaboration'],
    keywords: ['harmony', 'rhythm', 'resonance', 'emotion', 'pattern'],
    intensity: ['gentle', 'moderate', 'radical']
  },
  {
    id: 'quantum_physics',
    name: 'Quantum Physics',
    category: 'Science',
    description: 'Counterintuitive principles of reality at the smallest scales',
    examples: ['Quantum entanglement for instant connections', 'Superposition for multiple states', 'Observer effect on outcomes'],
    keywords: ['uncertainty', 'entanglement', 'superposition', 'probability', 'observation'],
    intensity: ['moderate', 'radical']
  },
  {
    id: 'cooking',
    name: 'Culinary Arts',
    category: 'Crafts',
    description: 'Flavor combinations, timing, temperature, and transformation processes',
    examples: ['Mise en place for project preparation', 'Flavor pairing for idea combinations', 'Fermentation for gradual development'],
    keywords: ['transformation', 'combination', 'timing', 'balance', 'craft'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'ecology',
    name: 'Ecosystem Dynamics',
    category: 'Nature',
    description: 'Complex interdependencies, feedback loops, and emergent behaviors',
    examples: ['Symbiotic partnerships', 'Succession patterns', 'Keystone species effects'],
    keywords: ['interdependence', 'cycles', 'balance', 'emergence', 'resilience'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'theater',
    name: 'Theater & Performance',
    category: 'Arts',
    description: 'Character development, dramatic structure, and audience engagement',
    examples: ['Three-act structure for project phases', 'Method acting for user empathy', 'Improvisation for adaptability'],
    keywords: ['character', 'narrative', 'audience', 'transformation', 'presence'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'martial_arts',
    name: 'Martial Arts',
    category: 'Philosophy',
    description: 'Balance, timing, energy redirection, and mental discipline',
    examples: ['Using opponent\'s force in negotiations', 'Flow state in difficult tasks', 'Defensive positioning'],
    keywords: ['balance', 'timing', 'discipline', 'flow', 'strategy'],
    intensity: ['moderate', 'radical']
  },
  {
    id: 'astronomy',
    name: 'Astronomy',
    category: 'Science',
    description: 'Scale, cycles, gravitational relationships, and cosmic perspectives',
    examples: ['Orbital mechanics for sustainable cycles', 'Gravitational wells for user retention', 'Star formation for growth'],
    keywords: ['scale', 'cycles', 'gravity', 'perspective', 'formation'],
    intensity: ['moderate', 'radical']
  },
  {
    id: 'anthropology',
    name: 'Cultural Anthropology',
    category: 'Human Systems',
    description: 'Cultural patterns, rituals, social structures, and human universals',
    examples: ['Ritual design for user onboarding', 'Gift economies for community building', 'Rite of passage for skill mastery'],
    keywords: ['culture', 'ritual', 'community', 'meaning', 'tradition'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'architecture',
    name: 'Architecture',
    category: 'Design',
    description: 'Space, flow, structure, and human experience of built environments',
    examples: ['Flow patterns for user journeys', 'Load-bearing structures for system design', 'Natural light for attention'],
    keywords: ['space', 'flow', 'structure', 'experience', 'foundation'],
    intensity: ['gentle', 'moderate']
  },
  {
    id: 'neuroscience',
    name: 'Neuroscience',
    category: 'Science',
    description: 'Brain networks, learning mechanisms, and cognitive biases',
    examples: ['Neural pathways for habit formation', 'Plasticity for adaptive systems', 'Mirror neurons for empathy'],
    keywords: ['networks', 'plasticity', 'learning', 'patterns', 'connection'],
    intensity: ['moderate', 'radical']
  },
  {
    id: 'mythology',
    name: 'World Mythology',
    category: 'Cultural',
    description: 'Universal stories, archetypal patterns, and symbolic meaning',
    examples: ['Hero\'s journey for user transformation', 'Creation myths for origin stories', 'Trickster figures for innovation'],
    keywords: ['archetype', 'journey', 'transformation', 'symbol', 'universal'],
    intensity: ['moderate', 'radical']
  },
  {
    id: 'economics',
    name: 'Economic Systems',
    category: 'Social Systems',
    description: 'Incentives, markets, scarcity, and value creation mechanisms',
    examples: ['Network effects for growth', 'Scarcity for desire', 'Liquidity for accessibility'],
    keywords: ['incentives', 'value', 'exchange', 'scarcity', 'networks'],
    intensity: ['gentle', 'moderate']
  }
];