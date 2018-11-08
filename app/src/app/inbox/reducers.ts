import { createFeatureSelector, createSelector, MetaReducer } from '@ngrx/store';
import { environment } from '../../environments/environment';
import { selectRecordEntities } from '../core/store';
import * as fromInbox from './store/inbox.reducer';

export const reducers = fromInbox.reducer;
export {State} from './store/inbox.reducer';

export const metaReducers: MetaReducer<fromInbox.State>[] = !environment.production ? [] : [];

export const selectFeature = createFeatureSelector<fromInbox.State>('inbox');
export const selectSelectedIds = createSelector(selectFeature, (state: fromInbox.State) => state.selectedIds);
export const selectSelectedRecords = createSelector(selectSelectedIds, selectRecordEntities, (ids, records) => ids.map(id => records[id]));
export const selectUnreadIds = createSelector(selectFeature, (state: fromInbox.State) => state.unreadIds);
export const selectUnreadRecords = createSelector(selectUnreadIds, selectRecordEntities);
export const selectMultiselect = createSelector(selectFeature, (state: fromInbox.State) => state.multiselect);
