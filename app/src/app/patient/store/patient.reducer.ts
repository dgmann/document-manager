import { PatientActions, PatientActionTypes } from './patient.actions';
import { Patient } from "./patient.model";

export interface State {
  selectedPatient: Patient;
  selectedRecordId: string;
  patientRecordMap: { [id: string]: string[] };
  filter: { categoryIds: string[], tags: string[] }
}

export const initialState: State = {
  selectedPatient: null,
  selectedRecordId: null,
  patientRecordMap: {},
  filter: {categoryIds: [], tags: []}
};

export function reducer(state = initialState, action: PatientActions): State {
  switch (action.type) {

    case PatientActionTypes.SetPatient:
      return Object.assign({}, state, {
        selectedPatient: action.payload.patient
      });
    case PatientActionTypes.SetPatientRecords:
      let map = Object.assign({}, state.patientRecordMap);
      map[action.payload.id] = action.payload.recordIds;
      return Object.assign({}, state, {
        patientRecordMap: map
      });
    case PatientActionTypes.SetFilter:
      return Object.assign({}, state, {
        filter: {categoryIds: action.payload.categoryIds, tags: action.payload.tags}
      });
    case PatientActionTypes.SelectRecord:
      return Object.assign({}, state, {
        selectedRecordId: action.payload.id
      });

    default:
      return state;
  }
}
