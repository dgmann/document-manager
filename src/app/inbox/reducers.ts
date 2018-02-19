import {createFeatureSelector, createSelector, MetaReducer} from '@ngrx/store';
import {environment} from '../../environments/environment';
import * as fromInbox from './store/inbox.reducer';

export const reducers = fromInbox.reducer;
export {State} from './store/inbox.reducer';

export const metaReducers: MetaReducer<fromInbox.State>[] = !environment.production ? [] : [];

export const selectFeature = createFeatureSelector<fromInbox.State>('inbox');
export const selectSelectedIds = createSelector(selectFeature, (state: fromInbox.State) => state.selectedIds);
