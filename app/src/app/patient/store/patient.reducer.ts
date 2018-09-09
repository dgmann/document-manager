import { PatientActions, PatientActionTypes } from './patient.actions';
import { Patient } from "./patient.model";
import { Moment } from "moment";

export interface Filter {
  categoryIds?: string[],
  tags?: string[],
  from?: Moment,
  until?: Moment
}

export interface State {
  selectedPatient: Patient;
  selectedRecordId: string;
  selectedCategory: string;
  patientRecordMap: { [id: string]: string[] };
  filter: Filter
}

export const initialState: State = {
  selectedPatient: null,
  selectedRecordId: null,
  selectedCategory: null,
  patientRecordMap: {},
  filter: {
    categoryIds: [],
    tags: [],
    from: null,
    until: null
  }
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
      let filter = state.filter;
      filter = Object.assign({}, filter, action.payload);
      return Object.assign({}, state, {
        filter: filter
      });
    case PatientActionTypes.SelectRecord:
      return Object.assign({}, state, {
        selectedRecordId: action.payload.id
      });

    case PatientActionTypes.SelectCategory:
      return {
        ...state,
        selectedCategory: action.payload.id
      };

    default:
      return state;
  }
}
