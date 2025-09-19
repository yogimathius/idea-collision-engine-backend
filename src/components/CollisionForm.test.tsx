import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { CollisionForm } from './CollisionForm';

describe('CollisionForm', () => {
  const mockOnSubmit = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render all form fields', () => {
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    expect(screen.getByLabelText(/your interests/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/current project/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/project type/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/collision intensity/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /generate idea collision/i })).toBeInTheDocument();
  });

  it('should show loading state when generating', () => {
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={true} />);

    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
    expect(screen.getByText(/generating collision/i)).toBeInTheDocument();
  });

  it('should require project description', async () => {
    const user = userEvent.setup();
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    const button = screen.getByRole('button', { name: /generate idea collision/i });
    expect(button).toBeDisabled();

    // Fill in project description
    const projectInput = screen.getByLabelText(/current project/i);
    await user.type(projectInput, 'Building a task manager');

    expect(button).toBeEnabled();
  });

  it('should submit form with correct data', async () => {
    const user = userEvent.setup();
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    // Fill in form
    await user.type(screen.getByLabelText(/your interests/i), 'productivity, design, technology');
    await user.type(screen.getByLabelText(/current project/i), 'Building a task management app');
    await user.selectOptions(screen.getByLabelText(/project type/i), 'product');
    await user.selectOptions(screen.getByLabelText(/collision intensity/i), 'moderate');

    // Submit form
    await user.click(screen.getByRole('button', { name: /generate idea collision/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith({
        userInterests: ['productivity', 'design', 'technology'],
        currentProject: 'Building a task management app',
        projectType: 'product',
        collisionIntensity: 'moderate'
      });
    });
  });

  it('should handle empty interests field', async () => {
    const user = userEvent.setup();
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    await user.type(screen.getByLabelText(/current project/i), 'Test project');
    await user.click(screen.getByRole('button', { name: /generate idea collision/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          userInterests: []
        })
      );
    });
  });

  it('should trim whitespace from interests and project', async () => {
    const user = userEvent.setup();
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    await user.type(screen.getByLabelText(/your interests/i), '  productivity , design  , technology  ');
    await user.type(screen.getByLabelText(/current project/i), '  Building a task manager  ');
    await user.click(screen.getByRole('button', { name: /generate idea collision/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          userInterests: ['productivity', 'design', 'technology'],
          currentProject: 'Building a task manager'
        })
      );
    });
  });

  it('should filter out empty interests', async () => {
    const user = userEvent.setup();
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    await user.type(screen.getByLabelText(/your interests/i), 'productivity,,design,  ,technology');
    await user.type(screen.getByLabelText(/current project/i), 'Test project');
    await user.click(screen.getByRole('button', { name: /generate idea collision/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          userInterests: ['productivity', 'design', 'technology']
        })
      );
    });
  });

  it('should prevent form submission when project is empty', async () => {
    // Mock window.alert
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {});
    
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    // Try to submit with empty project using the form element directly
    const form = screen.getByRole('button').closest('form');
    fireEvent.submit(form!);

    expect(alertSpy).toHaveBeenCalledWith('Please describe your current project');
    expect(mockOnSubmit).not.toHaveBeenCalled();

    alertSpy.mockRestore();
  });

  it('should have correct default values', () => {
    render(<CollisionForm onSubmit={mockOnSubmit} isGenerating={false} />);

    expect(screen.getByDisplayValue('Moderate (surprising but logical)')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Product')).toBeInTheDocument();
  });
});