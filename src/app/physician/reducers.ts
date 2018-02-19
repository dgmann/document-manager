import {createFeatureSelector, createSelector, MetaReducer} from '@ngrx/store';
import {environment} from '../../environments/environment';
import {selectRecordEntities} from "../store/record/record.selectors";
import * as fromPhysician from './store/physician.reducer';

export const reducers = fromPhysician.reducer;
export {State} from './store/physician.reducer';

export const metaReducers: MetaReducer<fromPhysician.State>[] = !environment.production ? [] : [];

export const selectFeature = createFeatureSelector<fromPhysician.State>('physician');
export const selectSelectedIds = createSelector(selectFeature, (state: fromPhysician.State) => state.selectedIds);
export const selectSelectedRecords = createSelector(selectSelectedIds, selectRecordEntities, (ids, records) => ids.map(id => records[id]));
