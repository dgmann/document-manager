import {createFeatureSelector, createSelector, MetaReducer} from '@ngrx/store';
import {environment} from '../../environments/environment';
import {Patient} from "../shared";
import {selectRecordEntities} from "../store/record/record.selectors";
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
