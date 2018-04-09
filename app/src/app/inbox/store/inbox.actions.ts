import {Action} from '@ngrx/store';

export enum InboxActionTypes {
  LoadRecords = '[Inbox] Load Records',
  SelectRecords = '[Inbox] Select Records',
  AddUnreadRecords = '[Inbox] Add unread Records',
  RemoveUnreadRecords = '[Inbox] Remove unread Records',
  SetMultiSelect = '[Inbox] Set Multiselect'
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

export class AddUnreadRecords implements Action {
  readonly type = InboxActionTypes.AddUnreadRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class RemoveUnreadRecords implements Action {
  readonly type = InboxActionTypes.RemoveUnreadRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class SetMultiSelect implements Action {
  readonly type = InboxActionTypes.SetMultiSelect;

  constructor(public payload: { multiselect: boolean }) {
  }
}

export type InboxActions =
  LoadRecords
  | SelectRecords
  | AddUnreadRecords
  | RemoveUnreadRecords
  | SetMultiSelect;
