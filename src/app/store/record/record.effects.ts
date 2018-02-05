import { Injectable } from '@angular/core';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { Observable } from "rxjs/Observable";
import { Action } from "@ngrx/store";
import {
  DeleteRecord,
  DeleteRecordFail,
  DeleteRecordSuccess,
  LoadRecordsFail,
  LoadRecordsSuccess,
  RecordActionTypes,
  UpdateRecord,
  UpdateRecordFail,
  UpdateRecordSuccess
} from "./record.actions";
import { catchError, map, mergeMap } from "rxjs/operators";
import { HttpClient } from "@angular/common/http";
import { Record } from "./record.model";
import { of } from "rxjs/observable/of";
import { environment } from "../../../environments/environment"


@Injectable()
export class RecordEffects {

  @Effect() load: Observable<Action> = this.actions$.pipe(
    ofType(RecordActionTypes.LoadRecords),
    mergeMap(action =>
      this.http.get<Record[]>(environment.api + '/records').pipe(
        map(data => new LoadRecordsSuccess({records: data})),
        catchError(err => of(new LoadRecordsFail({error: err})))
      )
    )
  );

  @Effect() update: Observable<Action> = this.actions$.pipe(
    ofType(RecordActionTypes.UpdateRecord),
    mergeMap((action: UpdateRecord) =>
      this.http.patch<Record>(`${environment.api}/records/${action.payload.record.id}`, action.payload.record.changes).pipe(
        map(data => new UpdateRecordSuccess({record: {id: action.payload.record.id as string, changes: data}})),
        catchError(err => of(new UpdateRecordFail({error: err})))
      )
    )
  );

  @Effect() delete: Observable<Action> = this.actions$.pipe(
    ofType(RecordActionTypes.DeleteRecord),
    mergeMap((action: DeleteRecord) =>
      this.http.delete(`${environment.api}/records/${action.payload.id}`).pipe(
        map(() => new DeleteRecordSuccess({id: action.payload.id})),
        catchError(err => of(new DeleteRecordFail({error: err})))
      )
    )
  );

  constructor(private actions$: Actions, private http: HttpClient) {
  }
}
