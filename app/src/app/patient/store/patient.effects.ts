import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Actions, createEffect, ofType} from '@ngrx/effects';
import {Action} from '@ngrx/store';
import {of} from 'rxjs';
import {catchError, map, switchMap} from 'rxjs/operators';
import {Record, LoadRecordsFail, LoadRecordsSuccess} from '@app/core/records';
import {PatientActionTypes, SelectPatient, SetPatient, SetPatientRecords} from './patient.actions';
import {Patient} from '@app/patient';
import {ConfigService} from '@app/core/config';

@Injectable()
export class PatientEffects {

  selectLoadPatientEffect$ = createEffect(() => this.actions$.pipe(
    ofType(PatientActionTypes.SelectPatientId),
    switchMap((action: SelectPatient) => this.http.get<Patient>(`${this.config.getApiUrl()}/patients/${action.payload.id}`)),
    map(data => new SetPatient({patient: data})),
    catchError(() => of(new SetPatient({patient: null})))
  ));

  selectLoadRecordsEffect$ = createEffect(() => this.actions$.pipe(
    ofType(PatientActionTypes.SelectPatientId),
    switchMap((action: SelectPatient) => this.http.get<Record[]>(`${this.config.getApiUrl()}/patients/${action.payload.id}/records`).pipe(
      switchMap(data => of<Action>(new LoadRecordsSuccess({records: data}), new SetPatientRecords({
        id: action.payload.id,
        recordIds: data.map(r => r.id)
      }))),
      catchError(err => of(new LoadRecordsFail({error: err})))
    ))
  ));

  constructor(private actions$: Actions, private http: HttpClient, private config: ConfigService) {
  }
}
