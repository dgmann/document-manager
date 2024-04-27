import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Actions, createEffect, ofType} from '@ngrx/effects';
import {Action} from '@ngrx/store';
import {Observable, of} from 'rxjs';
import {catchError, map, mergeMap} from 'rxjs/operators';
import {
  DeleteRecord,
  DeleteRecordFail,
  DeleteRecordSuccess,
  LoadRecords,
  LoadRecordsFail,
  LoadRecordsSuccess,
  RecordActionTypes,
  RecordErrorActions,
  UpdatePages,
  UpdateRecord,
  UpdateRecordFail,
  UpdateRecordSuccess
} from './record.actions';
import {Record} from './record.model';
import {AddError} from '../log/log.actions';
import {ConfigService} from '@app/core/config';


@Injectable()
export class RecordEffects {

  load: Observable<Action> = createEffect(() => this.actions$.pipe(
    ofType<LoadRecords>(RecordActionTypes.LoadRecords),
    mergeMap(action => {
      if (action.payload.query.id) {
        return this.http.get<Record>(this.config.getApiUrl() + '/records/' + action.payload.query.id).pipe(
          map(data => new LoadRecordsSuccess({records: [data]})),
          catchError(err => of(new LoadRecordsFail({error: err})))
        );
      } else {
        return this.http.get<Record[]>(this.config.getApiUrl() + '/records', {params: action.payload.query}).pipe(
          map(data => new LoadRecordsSuccess({records: data})),
          catchError(err => of(new LoadRecordsFail({error: err})))
        );
      }
    })
  ));

  update: Observable<Action> = createEffect(() => this.actions$.pipe(
    ofType<UpdateRecord>(RecordActionTypes.UpdateRecord),
    mergeMap(action =>
      this.http.patch<Record>(`${this.config.getApiUrl()}/records/${action.payload.record.id}`, action.payload.record.changes).pipe(
        map(data => new UpdateRecordSuccess({record: {id: action.payload.record.id as string, changes: data}})),
        catchError(err => of(new UpdateRecordFail({error: err})))
      )
    )
  ));

  delete: Observable<Action> = createEffect(() => this.actions$.pipe(
    ofType<DeleteRecord>(RecordActionTypes.DeleteRecord),
    mergeMap(action =>
      this.http.delete(`${this.config.getApiUrl()}/records/${action.payload.id}`).pipe(
        map(() => new DeleteRecordSuccess({id: action.payload.id})),
        catchError(err => of(new DeleteRecordFail({error: err})))
      )
    )
  ));

  updatePages: Observable<Action> = createEffect(() => this.actions$.pipe(
    ofType<UpdatePages>(RecordActionTypes.UpdatePages),
    mergeMap(action =>
      this.http.post<Record>(`${this.config.getApiUrl()}/records/${action.payload.id}/pages`, action.payload.updates).pipe(
        map(data => new UpdateRecordSuccess({record: {id: action.payload.id, changes: data}})),
        catchError(err => of(new UpdateRecordFail({error: err})))
      )
    )
  ));

  logError: Observable<Action> = createEffect(() => this.actions$.pipe(
    ofType<RecordErrorActions>(RecordActionTypes.LoadRecordsFail, RecordActionTypes.DeleteRecordFail, RecordActionTypes.UpdateRecordFail),
    mergeMap(action => of(this.errorFromHTTP(action.payload.error)))
  ));

  constructor(private actions$: Actions, private http: HttpClient, private config: ConfigService) {
  }

  private errorFromHTTP(err: HttpErrorResponse) {
    return new AddError({error: err.message});
  }
}
