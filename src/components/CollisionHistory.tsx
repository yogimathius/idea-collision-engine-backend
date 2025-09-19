import type { CollisionResult } from '../types';

interface CollisionHistoryProps {
  collisions: CollisionResult[];
  onSelectCollision?: (collision: CollisionResult) => void;
}

export function CollisionHistory({ collisions, onSelectCollision }: CollisionHistoryProps) {
  if (collisions.length === 0) {
    return (
      <div className="bg-gray-50 rounded-lg p-8 text-center">
        <div className="text-gray-400 text-6xl mb-4">🔍</div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">No collision history yet</h3>
        <p className="text-gray-600">Generate your first collision to start building your creative exploration history.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold text-gray-900">Recent Collisions</h2>
      <div className="grid gap-4">
        {collisions.slice(0, 10).map((collision, index) => (
          <div
            key={collision.id}
            onClick={() => onSelectCollision?.(collision)}
            className={`bg-white border rounded-lg p-4 ${
              onSelectCollision ? 'cursor-pointer hover:border-blue-300 hover:shadow-sm' : ''
            } transition-all`}
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="font-medium text-gray-900 text-sm mb-1">
                  {collision.primaryDomain} × {collision.collisionDomain}
                </div>
                <p className="text-xs text-gray-600 mb-2 line-clamp-2">
                  {collision.connection}
                </p>
                <div className="flex items-center space-x-4 text-xs text-gray-500">
                  <span>{collision.sparkQuestions.length} questions</span>
                  <span>Quality: {(collision.qualityScore * 100).toFixed(0)}%</span>
                  <span>{new Date(collision.timestamp).toLocaleDateString()}</span>
                </div>
              </div>
              {index === 0 && (
                <span className="bg-blue-100 text-blue-700 text-xs px-2 py-1 rounded-full ml-2">
                  Latest
                </span>
              )}
            </div>
          </div>
        ))}
      </div>
      
      {collisions.length > 10 && (
        <div className="text-center">
          <p className="text-sm text-gray-500">
            Showing 10 most recent collisions. Total: {collisions.length}
          </p>
        </div>
      )}
    </div>
  );
}