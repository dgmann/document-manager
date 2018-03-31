import {Injectable} from '@angular/core';
import {Actions, Effect, ofType} from '@ngrx/effects';
import {of} from "rxjs/observable/of";
import {map, mergeMap} from "rxjs/operators";
import {LoadRecordsSuccess, RecordActionTypes} from "../../store/record/record.actions";
import {AddUnreadRecords, InboxActionTypes, RemoveUnreadRecords, SelectRecords} from "./inbox.actions";

@Injectable()
export class InboxEffects {

  @Effect()
  addEffect$ = this.actions$.pipe(
    ofType(RecordActionTypes.LoadRecordsSuccess),
    map((action: LoadRecordsSuccess) => action.payload.records.map(record => record.id)),
    mergeMap(ids => of(new AddUnreadRecords({ids: ids})))
  );

  @Effect()
  removeEffect$ = this.actions$.pipe(
    ofType(InboxActionTypes.SelectRecords),
    mergeMap((action: SelectRecords) => of(new RemoveUnreadRecords({ids: action.payload.ids})))
  );

  constructor(private actions$: Actions) {
  }
}
