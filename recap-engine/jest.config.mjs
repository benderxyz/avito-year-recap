const transform = {
  '^.+\\.[jt]sx?$': [
    '@swc/jest',
    {
      jsc: {
        parser: {
          syntax: 'typescript',
          tsx: true,
        },
        transform: {
          react: {
            runtime: 'automatic',
          },
        },
      },
      module: {
        type: 'commonjs',
      },
    },
  ],
};

/** @type {import('jest').Config} */
export default {
  projects: [
    {
      displayName: 'core',
      rootDir: '<rootDir>/packages/core',
      testEnvironment: 'node',
      testMatch: ['<rootDir>/src/**/*.test.ts'],
      transform,
    },
    {
      displayName: 'react',
      rootDir: '<rootDir>/packages/react',
      testEnvironment: 'jsdom',
      testMatch: ['<rootDir>/src/**/*.test.{ts,tsx}'],
      transform,
      moduleNameMapper: {
        '^@recap-engine/core$': '<rootDir>/../core/src/index.ts',
      },
      setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
    },
  ],
  collectCoverageFrom: [
    'packages/core/src/**/*.ts',
    'packages/react/src/**/*.{ts,tsx}',
    '!**/*.test.{ts,tsx}',
    '!**/index.ts',
    '!**/interface.ts',
    '!packages/core/src/types/**',
    '!packages/react/src/test/**',
  ],
  coverageDirectory: '<rootDir>/coverage',
  coverageReporters: ['text', 'lcov', 'json-summary'],
  coverageThreshold: {
    global: {
      statements: 80,
      branches: 75,
      functions: 80,
      lines: 80,
    },
  },
};
