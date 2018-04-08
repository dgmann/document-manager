import {Injectable} from '@angular/core';
import {Actions, Effect, ofType} from '@ngrx/effects';
import {Action} from "@ngrx/store";
import {has} from 'lodash-es';
import {empty} from "rxjs/observable/empty";
import {from} from "rxjs/observable/from";
import {of} from "rxjs/observable/of";
import {map, switchMap} from "rxjs/operators";
import {RequiredAction} from "../../store";
import {
  DeleteRecordSuccess,
  LoadRecordsSuccess,
  RecordActionTypes,
  UpdateRecordSuccess
} from "../../store/record/record.actions";
import {AddRecord, RemoveRecord} from "./physician.actions";

@Injectable()
export class PhysicianEffects {

  @Effect()
  addEffect$ = this.actions$.pipe(
    ofType(RecordActionTypes.LoadRecordsSuccess),
    map((action: LoadRecordsSuccess) => action.payload.records),
    switchMap(records => of(...records.filter(record => record.requiredAction !== RequiredAction.NONE).map(record => new AddRecord({
      id: record.id,
      requiredAction: record.requiredAction
    }))))
  );

  @Effect()
  removeEffect$ = this.actions$.pipe(
    ofType(RecordActionTypes.DeleteRecordSuccess),
    map((action: DeleteRecordSuccess) => action.payload.id),
    switchMap(id => of(new RemoveRecord({id: id})))
  );

  @Effect()
  updateEffect$ = this.actions$.pipe(
    ofType(RecordActionTypes.UpdateRecordSuccess),
    map((action: UpdateRecordSuccess) => action.payload.record),
    switchMap(record => {
      if (has(record.changes, 'requiredAction')) {
        let actions: Action[] = [new RemoveRecord({id: record.id + ''})];
        if (record.changes.requiredAction !== RequiredAction.NONE) {
          actions.push(new AddRecord({id: record.id + '', requiredAction: record.changes.requiredAction}))
        }
        return from<Action>(actions);
      } else {
        return empty();
      }
    })
  );

  constructor(private actions$: Actions) {
  }
}
