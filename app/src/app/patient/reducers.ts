import {createFeatureSelector, createSelector, MetaReducer} from '@ngrx/store';
import {difference, flow, includes, intersection, reverse, sortBy} from 'lodash-es';
import {environment} from '@env/environment';
import {Record, selectDoneIds, selectRecordEntities} from '../core/records';
import * as fromPatient from './store/patient.reducer';
import {Filter} from './store/patient.reducer';
import { Dictionary } from '@ngrx/entity';

export const reducers = fromPatient.reducer;
export {State} from './store/patient.reducer';

export const metaReducers: MetaReducer<fromPatient.State>[] = !environment.production ? [] : [];

export const selectFeature = createFeatureSelector<fromPatient.State>('patient');
export const selectSelectedId = createSelector(selectFeature, state => state.selectedId);
export const selectSelectedPatient = createSelector(selectFeature, (state: fromPatient.State) => state.selectedPatient);
export const selectPatientRecordIds = createSelector(selectFeature, selectSelectedId, (state: fromPatient.State, id: string) => {
  if (!id) {
    return [];
  }
  return state.patientRecordMap[id] || [];
});
export const selectPatientRecords = createSelector(
  selectPatientRecordIds,
  selectDoneIds,
  selectRecordEntities,
  (ids: string[], doneIds: string[], records: Dictionary<Record>) => flow(
    ids => ids.map(id => records[id]) as Record[],
    records => sortBy(records, 'date'),
    records => reverse(records)
  )(intersection(ids, doneIds)));
export const selectFilter = createSelector(selectFeature, (state: fromPatient.State) => state.filter);
export const selectSelectedRecordId = createSelector(selectFeature, (state: fromPatient.State) => state.selectedRecordId);
export const selectSelectedRecord = createSelector(selectSelectedRecordId, selectRecordEntities, (id: string, records) => records[id]);

const filterRecords = (records: Record[], filter: Filter) => {
  return records.filter(record => {
    let result = false;
    if (filter.categoryIds.length === 0) {
      result = true;
    } else {
      result = includes(filter.categoryIds, record.category);
    }
    result = result && difference(filter.tags, record.tags).length === 0;

    if (filter.from) {
      result = result && record.date.isAfter(filter.from, 'day');
    }
    if (filter.until) {
      result = result && record.date.isBefore(filter.until, 'day');
    }

    return result;
  });
};
export const selectFilteredPatientRecords = createSelector(selectPatientRecords, selectFilter, filterRecords);

