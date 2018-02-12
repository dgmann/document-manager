import {Action} from '@ngrx/store';

export enum InboxActionTypes {
  AddRecords = '[Inbox] Add Records',
  RemoveRecords = '[Inbox] Remove Records',
  SelectRecords = '[Inbox] Select Records'
}

export class AddRecords implements Action {
  readonly type = InboxActionTypes.AddRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class RemoveRecords implements Action {
  readonly type = InboxActionTypes.RemoveRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class SelectRecords implements Action {
  readonly type = InboxActionTypes.SelectRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export type InboxActions =
  AddRecords
  | RemoveRecords
  | SelectRecords;
