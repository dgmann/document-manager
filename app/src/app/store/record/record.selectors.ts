import {Dictionary} from "@ngrx/entity/src/models";
import {createFeatureSelector, createSelector} from "@ngrx/store";
import {Record} from "./record.model";
import * as fromRecord from './record.reducer';

export const selectRecordState = createFeatureSelector<fromRecord.State>('records');

export const selectRecordIds = createSelector(selectRecordState, fromRecord.selectIds);
export const selectRecordEntities = createSelector(selectRecordState, fromRecord.selectEntities);
export const selectAllRecords = createSelector(selectRecordState, fromRecord.selectAll);
export const selectInvalidIds = createSelector(selectRecordState, (state: fromRecord.State) => state.invalidIds);

export const selectInboxIds = createSelector(selectRecordState, (state: fromRecord.State) => state.inboxIds);
export const selectEscalatedIds = createSelector(selectRecordState, (state: fromRecord.State) => state.escalatedIds);
export const selectReviewIds = createSelector(selectRecordState, (state: fromRecord.State) => state.reviewIds);
export const selectOtherIds = createSelector(selectRecordState, (state: fromRecord.State) => state.otherIds);
export const selectDoneIds = createSelector(selectRecordState, (state: fromRecord.State) => state.doneIds);

export const selectInboxRecords = createSelector(selectInboxIds, selectRecordEntities, (ids: string[], records: Dictionary<Record>) => ids.map(id => records[id]));
export const selectEscalatedRecords = createSelector(selectEscalatedIds, selectRecordEntities, (ids: string[], records: Dictionary<Record>) => ids.map(id => records[id]));
export const selectReviewRecords = createSelector(selectReviewIds, selectRecordEntities, (ids: string[], records: Dictionary<Record>) => ids.map(id => records[id]));
export const selectOtherRecords = createSelector(selectOtherIds, selectRecordEntities, (ids: string[], records: Dictionary<Record>) => ids.map(id => records[id]));
export const selectDoneIRecords = createSelector(selectDoneIds, selectRecordEntities, (ids: string[], records: Dictionary<Record>) => ids.map(id => records[id]));
