import { createFeatureSelector, createSelector, MetaReducer } from '@ngrx/store';
import { difference, includes } from 'lodash-es';
import { environment } from '../../environments/environment';
import { Record } from "../store";
import { selectRecordEntities } from "../store/record/record.selectors";
import { Patient } from "./store/patient.model";
import * as fromPatient from './store/patient.reducer';

export const reducers = fromPatient.reducer;
export {State} from './store/patient.reducer';

export const metaReducers: MetaReducer<fromPatient.State>[] = !environment.production ? [] : [];

export const selectFeature = createFeatureSelector<fromPatient.State>('patient');
export const selectSelectedPatient = createSelector(selectFeature, (state: fromPatient.State) => state.selectedPatient);
export const selectPatientRecordIds = createSelector(selectFeature, selectSelectedPatient, (state: fromPatient.State, patient: Patient) => {
  if (!patient) {
    return [];
  }
  return state.patientRecordMap[patient.id] || [];
});
export const selectPatientRecords = createSelector(selectPatientRecordIds, selectRecordEntities, (ids: string[], records) => ids.map(id => records[id]));
export const selectFilter = createSelector(selectFeature, (state: fromPatient.State) => state.filter);
export const selectFilteredPatientRecords = createSelector(selectPatientRecords, selectFilter, (records: Record[], filter) =>
  records.filter(record => (filter.categoryIds.length == 0 ? true : includes(filter.categoryIds, record.categoryId)) //Displays all records if no category is selected
    && difference(filter.tags, record.tags).length === 0));
export const selectSelectedRecordId = createSelector(selectFeature, (state: fromPatient.State) => state.selectedRecordId);
export const selectSelectedRecord = createSelector(selectSelectedRecordId, selectRecordEntities, (id: string, records) => records[id]);
