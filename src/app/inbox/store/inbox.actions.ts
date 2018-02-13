import { Action } from '@ngrx/store';

export enum InboxActionTypes {
  LoadRecords = '[Inbox] Load Records',
  SelectRecords = '[Inbox] Select Records'
}

export class LoadRecords implements Action {
  readonly type = InboxActionTypes.LoadRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class SelectRecords implements Action {
  readonly type = InboxActionTypes.SelectRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export type InboxActions =
  LoadRecords
  | SelectRecords;
