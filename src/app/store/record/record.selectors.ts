import { createFeatureSelector, createSelector } from "@ngrx/store";
import * as fromRecord from './record.reducer';

export const selectUserState = createFeatureSelector<fromRecord.State>('records');

export const selectRecordIds = createSelector(selectUserState, fromRecord.selectIds);
export const selectRecordEntities = createSelector(selectUserState, fromRecord.selectEntities);
export const selectAllRecords = createSelector(selectUserState, fromRecord.selectAll);
