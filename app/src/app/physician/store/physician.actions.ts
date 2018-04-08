import {Action} from '@ngrx/store';

export enum PhysicianActionTypes {
  LoadRecords = '[Physician] Load Records',
  SelectRecords = '[Physician] Select Records',
  AddRecord = '[Physician] Add Record',
  RemoveRecord = '[Physician] Remove Record',
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

export class AddRecord implements Action {
  readonly type = PhysicianActionTypes.AddRecord;

  constructor(public payload: { id: string, requiredAction: string }) {
  }
}

export class RemoveRecord implements Action {
  readonly type = PhysicianActionTypes.RemoveRecord;

  constructor(public payload: { id: string }) {
  }
}

export type PhysicianActions =
  LoadRecords
  | SelectRecords
  | AddRecord
  | RemoveRecord;
