import {Action} from '@ngrx/store';
import {Update} from '@ngrx/entity';
import {PageUpdate, Record} from './record.model';

export enum RecordActionTypes {
  LoadRecords = '[Record] Load Records',
  LoadRecordsFail = '[Record] Load Records Fail',
  LoadRecordsSuccess = '[Record] Load Records Success',
  UpdateRecord = '[Record] Update Record',
  UpdateRecordSuccess = '[Record] Update Records Success',
  UpdateRecordFail = '[Record] Update Records Fail',
  DeleteRecord = '[Record] Delete Record',
  DeleteRecordSuccess = '[Record] Delete Record Success',
  DeleteRecordFail = '[Record] Delete Records Fail',
  ClearRecords = '[Record] Clear Records',
  UpdatePages = '[Record] Update Pages'
}

export class LoadRecords implements Action {
  readonly type = RecordActionTypes.LoadRecords;

  constructor(public payload: { query: { [param: string]: string } }) {
  }
}

export class LoadRecordsSuccess implements Action {
  readonly type = RecordActionTypes.LoadRecordsSuccess;

  constructor(public payload: { records: Record[] }) {
    this.payload.records = this.payload.records.map(record => {
      record.date = record.date && new Date(record.date) || null;
      record.receivedAt = record.receivedAt && new Date(record.receivedAt) || null;
      return record;
    });
  }
}

export class LoadRecordsFail implements Action {
  readonly type = RecordActionTypes.LoadRecordsFail;

  constructor(public payload: { error: any }) {
  }
}

export class UpdateRecord implements Action {
  readonly type = RecordActionTypes.UpdateRecord;

  constructor(public payload: { record: Update<Record> }) {
  }
}

export class UpdateRecordSuccess implements Action {
  readonly type = RecordActionTypes.UpdateRecordSuccess;

  constructor(public payload: { record: Update<Record> }) {
    this.payload.record.changes.date = this.payload.record.changes.date
      && new Date(this.payload.record.changes.date)
      || null;
    this.payload.record.changes.receivedAt = this.payload.record.changes.receivedAt
      && new Date(this.payload.record.changes.receivedAt)
      || null;
  }
}

export class UpdateRecordFail implements Action {
  readonly type = RecordActionTypes.UpdateRecordFail;

  constructor(public payload: { error: any }) {
  }
}

export class DeleteRecord implements Action {
  readonly type = RecordActionTypes.DeleteRecord;

  constructor(public payload: { id: string }) {
  }
}

export class DeleteRecordSuccess implements Action {
  readonly type = RecordActionTypes.DeleteRecordSuccess;

  constructor(public payload: { id: string }) {
  }
}

export class DeleteRecordFail implements Action {
  readonly type = RecordActionTypes.DeleteRecordFail;

  constructor(public payload: { error: any }) {
  }
}

export class ClearRecords implements Action {
  readonly type = RecordActionTypes.ClearRecords;

  constructor(public payload: {}) {
  }
}

export class UpdatePages implements Action {
  readonly type = RecordActionTypes.UpdatePages;

  constructor(public payload: { id: string, updates: PageUpdate[] }) {
  }
}

export type RecordActions =
  LoadRecords
  | LoadRecordsSuccess
  | LoadRecordsFail
  | UpdateRecord
  | UpdateRecordSuccess
  | UpdateRecordFail
  | DeleteRecord
  | DeleteRecordSuccess
  | DeleteRecordFail
  | ClearRecords
  | UpdatePages;

export type RecordErrorActions = LoadRecordsFail | UpdateRecordFail | DeleteRecordFail;
