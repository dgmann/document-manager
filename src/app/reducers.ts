import { ActionReducerMap, MetaReducer } from "@ngrx/store";
import { storeFreeze } from 'ngrx-store-freeze';

import { environment } from '../environments/environment';

/**
 * As mentioned, we treat each reducer like a table in a database. This means
 * our top level state interface is just a map of keys to inner state types.
 */
export interface State {}

export const reducers: ActionReducerMap<State> = {};

/**
 * By default, @ngrx/store uses combineReducers with the reducer map to compose
 * the root meta-reducer. To add more meta-reducers, provide an array of meta-reducers
 * that will be composed to form the root meta-reducer.
 */

export const metaReducers: MetaReducer<State>[] = !environment.production
  ? [storeFreeze]
  : [];
