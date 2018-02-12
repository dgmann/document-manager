import {Injectable} from '@angular/core';
import {Actions, Effect} from '@ngrx/effects';
import {InboxActionTypes} from './inbox.actions';

@Injectable()
export class InboxEffects {

  @Effect()
  effect$ = this.actions$.ofType(InboxActionTypes.SelectRecords);

  constructor(private actions$: Actions) {
  }
}
