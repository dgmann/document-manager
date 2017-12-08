import { createFeatureSelector, createSelector } from "@ngrx/store";
import * as fromRecord from "./record.reducer";

export const selectRecordState = createFeatureSelector<fromRecord.State>('records');

export const {
  // select the array of user ids
  selectIds: selectRecordIds,

  // select the dictionary of user entities
  selectEntities: selectRecordEntities,

  // select the array of users
  selectAll: selectAllRecords,

  // select the total user count
  selectTotal: selectRecordTotal
} = fromRecord.adapter.getSelectors(selectRecordState);


export const selectRecord = (id: string) => createSelector(selectRecordEntities, recordEntities => recordEntities[id]);

export const selectCurrentRecordId = createSelector(selectRecordState, fromRecord.getSelectedRecordId);
export const selectCurrentRecord = createSelector(
  selectRecordEntities,
  selectCurrentRecordId,
  (recordEntities, recordId) => recordEntities[recordId]
);
