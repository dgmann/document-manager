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
export const selectEscalatedIds = createSelector(selectFeature, (state: fromPhysician.State) => state.escalatedIds);
export const selectReviewIds = createSelector(selectFeature, (state: fromPhysician.State) => state.reviewIds);
export const selectOtherIds = createSelector(selectFeature, (state: fromPhysician.State) => state.otherIds);
export const selectEscalatedRecords = createSelector(selectEscalatedIds, selectRecordEntities, (ids, records) => ids.map(id => records[id]));
export const selectReviewRecords = createSelector(selectReviewIds, selectRecordEntities, (ids, records) => ids.map(id => records[id]));
export const selectOtherRecords = createSelector(selectOtherIds, selectRecordEntities, (ids, records) => ids.map(id => records[id]));
