import { Action } from '@ngrx/store';
import { Patient } from "./patient.model";

export enum PatientActionTypes {
  SelectPatientId = '[Patient] Select Patient ID',
  SetPatient = '[Patient] Set Patient',
  SetPatientRecords = '[Patient] Set Patient Record ids',
  SetFilter = '[Patient] Set Filter',
  SelectRecord = '[Patient] Select Record'
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

  constructor(public payload: { categoryIds: string[], tags: string[] }) {
  }
}

export class SelectRecord implements Action {
  readonly type = PatientActionTypes.SelectRecord;

  constructor(public payload: { id: string }) {
  }
}

export type PatientActions = SelectPatient
  | SetPatientRecords
  | SetPatient
  | SetFilter
  | SelectRecord;
