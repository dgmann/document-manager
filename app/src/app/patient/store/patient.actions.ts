import {Action} from '@ngrx/store';
import {Patient} from './patient.model';
import {Filter} from './patient.reducer';

export enum PatientActionTypes {
  SelectPatientId = '[Patient] Select Patient ID',
  SetPatient = '[Patient] Set Patient',
  SetPatientRecords = '[Patient] Set Patient Record ids',
  SetFilter = '[Patient] Set Filter',
  SelectRecord = '[Patient] Select Record',
  SelectCategory = '[Patient] Select Category'
}

export class SelectPatient implements Action {
  readonly type = PatientActionTypes.SelectPatientId;

  constructor(public payload: { id: string }) {
  }
}

export class SetPatient implements Action {
  readonly type = PatientActionTypes.SetPatient;

  constructor(public payload: { patient: Patient }) {
  }
}

export class SetPatientRecords implements Action {
  readonly type = PatientActionTypes.SetPatientRecords;

  constructor(public payload: { id: string, recordIds: string[] }) {
  }
}

export class SetFilter implements Action {
  readonly type = PatientActionTypes.SetFilter;

  constructor(public payload: Filter) {
  }
}

export class SelectRecord implements Action {
  readonly type = PatientActionTypes.SelectRecord;

  constructor(public payload: { id: string }) {
  }
}

export class SelectCategory implements Action {
  readonly type = PatientActionTypes.SelectCategory;

  constructor(public payload: { id: string }) {
  }
}

export type PatientActions = SelectPatient
  | SetPatientRecords
  | SetPatient
  | SetFilter
  | SelectRecord
  | SelectCategory;
