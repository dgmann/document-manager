import { ResourceDefinition } from 'ngrx-json-api';

export const resourceDefinitions: Array<ResourceDefinition> = [
  { type: 'Record', collectionPath: 'records' },
  { type: 'Person', collectionPath: 'people' },
];
