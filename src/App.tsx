import { useState, useEffect } from 'react';
import { CollisionForm } from './components/CollisionForm';
import { CollisionResult } from './components/CollisionResult';
import { CollisionHistory } from './components/CollisionHistory';
import { CollisionEngine } from './collisionEngine';
import type { CollisionInput, CollisionResult as CollisionResultType } from './types';
import { 
  saveCollisionToHistory, 
  getCollisionHistory, 
  updateCollisionRating, 
  addCollisionNotes 
} from './utils/localStorage';

function App() {
  const [currentCollision, setCurrentCollision] = useState<CollisionResultType | null>(null);
  const [collisionHistory, setCollisionHistory] = useState<CollisionResultType[]>([]);
  const [isGenerating, setIsGenerating] = useState(false);
  const [activeTab, setActiveTab] = useState<'generate' | 'history'>('generate');
  const [collisionEngine] = useState(() => new CollisionEngine());

  useEffect(() => {
    // Load collision history on app start
    const history = getCollisionHistory();
    setCollisionHistory(history);
  }, []);

  const handleGenerateCollision = async (input: CollisionInput) => {
    setIsGenerating(true);
    try {
      const collision = await collisionEngine.generateCollision(input);
      setCurrentCollision(collision);
      
      // Save to history
      saveCollisionToHistory(collision);
      setCollisionHistory(prev => [collision, ...prev]);
      
      // Auto-switch to show result
      setActiveTab('generate');
    } catch (error) {
      console.error('Failed to generate collision:', error);
      alert('Failed to generate collision. Please try again.');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleRateCollision = (rating: number) => {
    if (currentCollision) {
      updateCollisionRating(currentCollision.id, rating);
      setCurrentCollision(prev => prev ? { ...prev, rating } : null);
    }
  };

  const handleAddNotes = (notes: string) => {
    if (currentCollision) {
      addCollisionNotes(currentCollision.id, notes);
      setCurrentCollision(prev => prev ? { ...prev, notes } : null);
    }
  };

  const handleSelectHistoryCollision = (collision: CollisionResultType) => {
    setCurrentCollision(collision);
    setActiveTab('generate');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {/* Navigation */}
        <div className="flex justify-center mb-8">
          <div className="bg-white p-1 rounded-lg shadow-sm border">
            <button
              onClick={() => setActiveTab('generate')}
              className={`px-6 py-2 rounded-md text-sm font-medium transition-colors ${
                activeTab === 'generate'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Generate Collision
            </button>
            <button
              onClick={() => setActiveTab('history')}
              className={`px-6 py-2 rounded-md text-sm font-medium transition-colors ${
                activeTab === 'history'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              History ({collisionHistory.length})
            </button>
          </div>
        </div>

        {/* Content */}
        {activeTab === 'generate' ? (
          <div className="space-y-8">
            {/* Form */}
            <CollisionForm 
              onSubmit={handleGenerateCollision}
              isGenerating={isGenerating}
            />

            {/* Current Collision Result */}
            {currentCollision && (
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-xl font-semibold text-gray-900">Your Collision</h2>
                  <button
                    onClick={() => setCurrentCollision(null)}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    ✕
                  </button>
                </div>
                <CollisionResult
                  collision={currentCollision}
                  onRate={handleRateCollision}
                  onAddNotes={handleAddNotes}
                />
              </div>
            )}

            {/* Quick stats */}
            {collisionHistory.length > 0 && (
              <div className="bg-blue-50 rounded-lg p-4 text-center">
                <p className="text-blue-700">
                  <span className="font-semibold">{collisionHistory.length}</span> collisions generated
                  {collisionHistory.length >= 5 && ' • You\'re building creative momentum! 🚀'}
                </p>
              </div>
            )}
          </div>
        ) : (
          <CollisionHistory
            collisions={collisionHistory}
            onSelectCollision={handleSelectHistoryCollision}
          />
        )}

        {/* Footer */}
        <div className="mt-16 text-center text-gray-500 text-sm">
          <p>Idea Collision Engine • Break out of familiar patterns</p>
          <p className="mt-1">
            Generate unexpected but relevant connections to spark creative breakthroughs
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;