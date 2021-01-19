import {PatientActions, PatientActionTypes} from './patient.actions';
import {Patient} from './patient.model';
import {Moment} from 'moment';

export interface Filter {
  categoryIds?: string[];
  tags?: string[];
  from?: Moment;
  until?: Moment;
}

export interface State {
  selectedId: string;
  selectedPatient: Patient;
  selectedRecordId: string;
  patientRecordMap: { [id: string]: string[] };
  filter: Filter;
}

export const initialState: State = {
  selectedId: null,
  selectedPatient: null,
  selectedRecordId: null,
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

    case PatientActionTypes.SelectPatientId:
      return {
        ...state,
        selectedId: action.payload.id
      };

    case PatientActionTypes.SetPatient:
      return Object.assign({}, state, {
        selectedPatient: action.payload.patient
      });
    case PatientActionTypes.SetPatientRecords:
      const map = Object.assign({}, state.patientRecordMap);
      map[action.payload.id] = action.payload.recordIds;
      return Object.assign({}, state, {
        patientRecordMap: map
      });
    case PatientActionTypes.SetFilter:
      let filter = state.filter;
      filter = Object.assign({}, filter, action.payload);
      return Object.assign({}, state, {
        filter
      });
    case PatientActionTypes.SelectRecord:
      return Object.assign({}, state, {
        selectedRecordId: action.payload.id
      });

    default:
      return state;
  }
}
