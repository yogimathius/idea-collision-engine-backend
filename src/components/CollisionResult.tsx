import { useState } from 'react';
import type { CollisionResult as CollisionResultType } from '../types';

interface CollisionResultProps {
  collision: CollisionResultType;
  onRate?: (rating: number) => void;
  onAddNotes?: (notes: string) => void;
}

export function CollisionResult({ collision, onRate, onAddNotes }: CollisionResultProps) {
  const [rating, setRating] = useState<number | null>(null);
  const [notes, setNotes] = useState('');
  const [showNotes, setShowNotes] = useState(false);

  const handleRate = (newRating: number) => {
    setRating(newRating);
    onRate?.(newRating);
  };

  const handleSaveNotes = () => {
    onAddNotes?.(notes);
    setShowNotes(false);
  };

  return (
    <div className="bg-white rounded-xl shadow-sm border p-8 space-y-6 animate-in fade-in slide-in-from-top-4 duration-500">
      {/* Header */}
      <div className="text-center border-b pb-6">
        <div className="text-sm text-gray-500 mb-2">Collision between</div>
        <div className="text-2xl font-bold text-gray-900">
          {collision.primaryDomain} <span className="text-blue-500">×</span> {collision.collisionDomain}
        </div>
        <div className="text-sm text-gray-400 mt-2">
          Quality Score: {(collision.qualityScore * 100).toFixed(0)}%
        </div>
      </div>

      {/* Connection */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">🎯 The Connection</h3>
        <div className="bg-blue-50 p-4 rounded-lg">
          <p className="text-gray-700 leading-relaxed">{collision.connection}</p>
        </div>
      </div>

      {/* Spark Questions */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">💡 Spark Questions</h3>
        <ul className="space-y-2">
          {collision.sparkQuestions.map((question, index) => (
            <li key={index} className="flex items-start space-x-3">
              <span className="text-blue-500 font-semibold mt-1">Q{index + 1}:</span>
              <span className="text-gray-700">{question}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Examples */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">🌟 Real-World Examples</h3>
        <ul className="space-y-2">
          {collision.examples.map((example, index) => (
            <li key={index} className="flex items-start space-x-3">
              <span className="text-green-500 font-semibold">•</span>
              <span className="text-gray-700">{example}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Next Steps */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">🚀 Next Steps</h3>
        <ol className="space-y-2">
          {collision.nextSteps.map((step, index) => (
            <li key={index} className="flex items-start space-x-3">
              <span className="bg-purple-100 text-purple-700 rounded-full w-6 h-6 flex items-center justify-center text-sm font-semibold flex-shrink-0 mt-0.5">
                {index + 1}
              </span>
              <span className="text-gray-700">{step}</span>
            </li>
          ))}
        </ol>
      </div>

      {/* Rating */}
      <div className="border-t pt-6">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-gray-900">Rate this collision</h3>
            <p className="text-sm text-gray-500">How useful was this collision for your project?</p>
          </div>
          <div className="flex space-x-1">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                onClick={() => handleRate(star)}
                className={`text-2xl transition-colors ${
                  rating && star <= rating
                    ? 'text-yellow-400'
                    : 'text-gray-300 hover:text-yellow-200'
                }`}
              >
                ★
              </button>
            ))}
          </div>
        </div>

        {/* Notes */}
        <div className="mt-4">
          {!showNotes ? (
            <button
              onClick={() => setShowNotes(true)}
              className="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              + Add exploration notes
            </button>
          ) : (
            <div className="space-y-2">
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="What insights did this collision spark? What will you try next?"
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              />
              <div className="flex space-x-2">
                <button
                  onClick={handleSaveNotes}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md text-sm hover:bg-blue-700"
                >
                  Save Notes
                </button>
                <button
                  onClick={() => setShowNotes(false)}
                  className="px-4 py-2 text-gray-600 text-sm hover:text-gray-700"
                >
                  Cancel
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Timestamp */}
      <div className="text-xs text-gray-400 text-center">
        Generated at {collision.timestamp.toLocaleString()}
      </div>
    </div>
  );
}