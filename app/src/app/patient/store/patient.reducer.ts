import {Patient} from "../../shared";
import {PatientActions, PatientActionTypes} from './patient.actions';

export interface State {
  selectedPatient: Patient;
  patientRecordMap: { [id: string]: string[] };
}

export const initialState: State = {
  selectedPatient: null,
  patientRecordMap: {}
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


    default:
      return state;
  }
}
