/**
 * For a detailed explanation regarding each configuration property, visit:
 * https://jestjs.io/docs/configuration
 */

import type {Config} from 'jest';

const config: Config = {
  preset: 'jest-preset-angular',
  setupFilesAfterEnv: ['<rootDir>/setup-jest.ts'],
  transform: {
    '^.+\\.(ts|mjs|js|html)$': [
      'jest-preset-angular',
      {
        tsconfig: '<rootDir>/tsconfig.spec.json',
        stringifyContentPathRegex: '\\.(html|svg)$',
      },
    ],
  },
  transformIgnorePatterns: [
    'node_modules/(?!.*\\.mjs$)',
    'node_modules/lodash-es/*',
  ],

  moduleFileExtensions: ['ts', 'html', 'js', 'json'],
  moduleDirectories: ['node_modules', 'src'],
  modulePaths: ['<rootDir>'],
  moduleNameMapper: {
      '^@app/(.*)$': ['<rootDir>/src/app/$1'],
      '^@env/(.*)$': ['<rootDir>/src/environments/$1'],
      '^lodash-es$': 'lodash', '^lodash-es/(.*)': 'lodash/$1' 
  }
};

export default config;
