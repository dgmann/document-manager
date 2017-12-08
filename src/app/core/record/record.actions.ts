import { Record } from './record.model';
import { PayloadAction } from "../payload-action";

export const LOAD_RECORDS = '[Record] Load Records';
export const ADD_RECORD = '[Record] Add Record';
export const ADD_RECORDS = '[Record] Add Records';
export const UPDATE_RECORD = '[Record] Update Record';
export const UPDATE_RECORDS = '[Record] Update Records';
export const DELETE_RECORD = '[Record] Delete Record';
export const DELETE_RECORDS = '[Record] Delete Records';
export const CLEAR_RECORDS = '[Record] Clear Records';

export class LoadRecords implements PayloadAction {
  readonly type = LOAD_RECORDS;

  constructor(public payload: { records: Record[] }) {}
}

export class AddRecord implements PayloadAction {
  readonly type = ADD_RECORD;

  constructor(public payload: { records: Record }) {}
}

export class AddRecords implements PayloadAction {
  readonly type = ADD_RECORDS;

  constructor(public payload: { records: Record[] }) {}
}

export class UpdateRecord implements PayloadAction {
  readonly type = UPDATE_RECORD;

  constructor(public payload: { records: { id: string, changes: Record } }) {}
}

export class UpdateRecords implements PayloadAction {
  readonly type = UPDATE_RECORDS;

  constructor(public payload: { records: { id: string, changes: Record }[] }) {}
}

export class DeleteRecord implements PayloadAction {
  readonly type = DELETE_RECORD;

  constructor(public payload: { id: string }) {}
}

export class DeleteRecords implements PayloadAction {
  readonly type = DELETE_RECORDS;

  constructor(public payload: { ids: string[] }) {}
}

export class ClearRecords implements PayloadAction {
  readonly type = CLEAR_RECORDS;
  constructor(public payload: null) {}
}

export type All =
  LoadRecords
  | AddRecord
  | AddRecords
  | UpdateRecord
  | UpdateRecords
  | DeleteRecord
  | DeleteRecords
  | ClearRecords;
