import {createFeatureSelector, createSelector} from "@ngrx/store";
import * as fromRecord from './record.reducer';

export const selectRecordState = createFeatureSelector<fromRecord.State>('records');

export const selectRecordIds = createSelector(selectRecordState, fromRecord.selectIds);
export const selectRecordEntities = createSelector(selectRecordState, fromRecord.selectEntities);
export const selectAllRecords = createSelector(selectRecordState, fromRecord.selectAll);
