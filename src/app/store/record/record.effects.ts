import { Injectable } from '@angular/core';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { Observable } from "rxjs/Observable";
import { Action } from "@ngrx/store";
import { AddRecord, AddRecords, RecordActionTypes, UpdateRecord } from "./record.actions";
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
        map(data => new AddRecords({records: data})),
        catchError(() => of({type: 'LOAD_FAILED'}))
      )
    )
  );
  @Effect() update: Observable<Action> = this.actions$.pipe(
    ofType(RecordActionTypes.UpdateRecord),
    mergeMap((action: UpdateRecord) =>
      this.http.patch<Record>(`${environment.api}/records/${action.payload.record.id}`, action.payload.record.changes).pipe(
        map(data => new AddRecord({record: data})),
        catchError(() => of({type: 'UPDATE_FAILED'}))
      )
    )
  );

  constructor(private actions$: Actions, private http: HttpClient) {
  }
}
