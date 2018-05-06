import {HttpClient} from "@angular/common/http";
import {Injectable} from '@angular/core';
import {Actions, Effect, ofType} from '@ngrx/effects';
import {Action} from "@ngrx/store";
import {of} from "rxjs";
import {catchError, map, switchMap} from "rxjs/operators";
import {environment} from "../../../environments/environment";
import {Record} from "../../store";
import {LoadRecordsFail, LoadRecordsSuccess} from "../../store/record/record.actions";
import {PatientActionTypes, SelectPatient, SetPatient, SetPatientRecords} from './patient.actions';
import {Patient} from "./patient.model";

@Injectable()
export class PatientEffects {

  @Effect()
  selectLoadPatientEffect$ = this.actions$.pipe(
    ofType(PatientActionTypes.SelectPatientId),
    switchMap((action: SelectPatient) => this.http.get<Patient>(`${environment.api}/patients/${action.payload.id}`)),
    map(data => new SetPatient({patient: data})),
    catchError(_ => of(new SetPatient({patient: null})))
  );

  @Effect()
  selectLoadRecordsEffect$ = this.actions$.pipe(
    ofType(PatientActionTypes.SelectPatientId),
    switchMap((action: SelectPatient) => this.http.get<Record[]>(`${environment.api}/patients/${action.payload.id}/records`).pipe(
      switchMap(data => of<Action>(new LoadRecordsSuccess({records: data}), new SetPatientRecords({
        id: action.payload.id,
        recordIds: data.map(r => r.id)
      }))),
      catchError(err => of(new LoadRecordsFail({error: err})))
    ))
  );

  constructor(private actions$: Actions, private http: HttpClient) {
  }
}
