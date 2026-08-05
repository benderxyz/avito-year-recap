import '@testing-library/jest-dom';
import { afterEach, jest } from '@jest/globals';
import { cleanup } from '@testing-library/react';

afterEach(() => {
  cleanup();
  jest.restoreAllMocks();
});

Object.defineProperty(window, 'matchMedia', {
  configurable: true,
  writable: true,
  value: jest.fn().mockImplementation((query: unknown) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    addListener: jest.fn(),
    removeListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

if (!window.PointerEvent) {
  Object.defineProperty(window, 'PointerEvent', {
    configurable: true,
    writable: true,
    value: MouseEvent,
  });
}

if (!window.requestAnimationFrame) {
  window.requestAnimationFrame = (callback) =>
    window.setTimeout(() => callback(performance.now()), 16);
  window.cancelAnimationFrame = (id) => window.clearTimeout(id);
}
