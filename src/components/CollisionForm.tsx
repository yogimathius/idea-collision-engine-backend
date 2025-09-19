import { useState } from 'react';
import type { CollisionInput } from '../types';

interface CollisionFormProps {
  onSubmit: (input: CollisionInput) => void;
  isGenerating: boolean;
}

export function CollisionForm({ onSubmit, isGenerating }: CollisionFormProps) {
  const [interests, setInterests] = useState<string>('');
  const [currentProject, setCurrentProject] = useState('');
  const [projectType, setProjectType] = useState<CollisionInput['projectType']>('product');
  const [collisionIntensity, setCollisionIntensity] = useState<CollisionInput['collisionIntensity']>('moderate');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!currentProject.trim()) {
      alert('Please describe your current project');
      return;
    }

    const userInterests = interests.split(',').map(s => s.trim()).filter(s => s.length > 0);
    
    onSubmit({
      userInterests,
      currentProject: currentProject.trim(),
      projectType,
      collisionIntensity
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-8 rounded-xl shadow-sm border">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Idea Collision Engine</h1>
        <p className="text-gray-600">Break out of familiar patterns with unexpected but relevant idea combinations</p>
      </div>

      <div className="space-y-4">
        <div>
          <label htmlFor="interests" className="block text-sm font-medium text-gray-700 mb-2">
            Your Interests <span className="text-gray-400">(comma-separated)</span>
          </label>
          <input
            type="text"
            id="interests"
            value={interests}
            onChange={(e) => setInterests(e.target.value)}
            placeholder="productivity, meditation, javascript, design..."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <p className="text-xs text-gray-500 mt-1">What areas do you usually think about?</p>
        </div>

        <div>
          <label htmlFor="project" className="block text-sm font-medium text-gray-700 mb-2">
            Current Project <span className="text-red-400">*</span>
          </label>
          <textarea
            id="project"
            value={currentProject}
            onChange={(e) => setCurrentProject(e.target.value)}
            placeholder="Building a task management app for creative teams..."
            rows={3}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="projectType" className="block text-sm font-medium text-gray-700 mb-2">
              Project Type
            </label>
            <select
              id="projectType"
              value={projectType}
              onChange={(e) => setProjectType(e.target.value as CollisionInput['projectType'])}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="product">Product</option>
              <option value="content">Content</option>
              <option value="business">Business</option>
              <option value="research">Research</option>
            </select>
          </div>

          <div>
            <label htmlFor="intensity" className="block text-sm font-medium text-gray-700 mb-2">
              Collision Intensity
            </label>
            <select
              id="intensity"
              value={collisionIntensity}
              onChange={(e) => setCollisionIntensity(e.target.value as CollisionInput['collisionIntensity'])}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="gentle">Gentle (familiar connections)</option>
              <option value="moderate">Moderate (surprising but logical)</option>
              <option value="radical">Radical (mind-bending leaps)</option>
            </select>
          </div>
        </div>
      </div>

      <button
        type="submit"
        disabled={isGenerating || !currentProject.trim()}
        className="w-full bg-blue-600 text-white py-4 px-6 rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {isGenerating ? (
          <>
            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Generating Collision...
          </>
        ) : (
          'Generate Idea Collision 💥'
        )}
      </button>
    </form>
  );
}