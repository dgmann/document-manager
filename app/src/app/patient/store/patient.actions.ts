import {Action} from '@ngrx/store';
import {Patient} from "../../shared";

export enum PatientActionTypes {
  SelectPatientId = '[Patient] Select Patient ID',
  SetPatient = '[Patient] Set Patient',
  SetPatientRecords = '[Patient] Set Patient Record ids'
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

export type PatientActions = SelectPatient | SetPatientRecords | SetPatient;
