import { Action } from '@ngrx/store';
import { Update } from '@ngrx/entity';
import { Record } from './record.model';

export enum RecordActionTypes {
  LoadRecords = '[Record] Load Records',
  AddRecord = '[Record] Add Record',
  AddRecords = '[Record] Add Records',
  UpdateRecord = '[Record] Update Record',
  UpdateRecords = '[Record] Update Records',
  DeleteRecord = '[Record] Delete Record',
  DeleteRecords = '[Record] Delete Records',
  ClearRecords = '[Record] Clear Records'
}

export class LoadRecords implements Action {
  readonly type = RecordActionTypes.LoadRecords;

  constructor(public payload: { records: Record[] }) {
  }
}

export class AddRecord implements Action {
  readonly type = RecordActionTypes.AddRecord;

  constructor(public payload: { record: Record }) {
  }
}

export class AddRecords implements Action {
  readonly type = RecordActionTypes.AddRecords;

  constructor(public payload: { records: Record[] }) {
  }
}

export class UpdateRecord implements Action {
  readonly type = RecordActionTypes.UpdateRecord;

  constructor(public payload: { record: Update<Record> }) {
  }
}

export class UpdateRecords implements Action {
  readonly type = RecordActionTypes.UpdateRecords;

  constructor(public payload: { records: Update<Record>[] }) {
  }
}

export class DeleteRecord implements Action {
  readonly type = RecordActionTypes.DeleteRecord;

  constructor(public payload: { id: string }) {
  }
}

export class DeleteRecords implements Action {
  readonly type = RecordActionTypes.DeleteRecords;

  constructor(public payload: { ids: string[] }) {
  }
}

export class ClearRecords implements Action {
  readonly type = RecordActionTypes.ClearRecords;
}

export type RecordActions =
  LoadRecords
  | AddRecord
  | AddRecords
  | UpdateRecord
  | UpdateRecords
  | DeleteRecord
  | DeleteRecords
  | ClearRecords;
