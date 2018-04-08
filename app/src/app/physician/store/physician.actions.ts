import {Action} from '@ngrx/store';

export enum PhysicianActionTypes {
  LoadRecords = '[Physician] Load Records',
  SelectRecords = '[Physician] Select Records',
  SetRecord = '[Physician] Set Record',
}

export class LoadRecords implements Action {
  readonly type = PhysicianActionTypes.LoadRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class SelectRecords implements Action {
  readonly type = PhysicianActionTypes.SelectRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class SetRecord implements Action {
  readonly type = PhysicianActionTypes.SetRecord;

  constructor(public payload: { id: string, requiredAction: string }) {
  }
}

export type PhysicianActions =
  LoadRecords
  | SelectRecords
  | SetRecord
