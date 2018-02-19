import {Action} from '@ngrx/store';

export enum PhysicianActionTypes {
  LoadRecords = '[Physician] Load Records',
  SelectRecords = '[Physician] Select Records'
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

export type PhysicianActions =
  LoadRecords
  | SelectRecords;
