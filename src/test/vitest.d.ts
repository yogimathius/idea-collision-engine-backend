import type { TestingLibraryMatchers } from '@testing-library/jest-dom/matchers';

declare module 'vitest' {
  interface Assertion<T = unknown> extends TestingLibraryMatchers<T, void> {
    // Custom matchers can be added here
    toBeInTheDocument(): void;
    toBeEnabled(): void;  
    toBeDisabled(): void;
  }
  interface AsymmetricMatchersContaining<T = unknown> extends TestingLibraryMatchers<T, void> {
    // Custom matchers for asymmetric matching
    toBeInTheDocument(): void;
    toBeEnabled(): void;
    toBeDisabled(): void;
  }
}