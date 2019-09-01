import {Action} from '@ngrx/store';
import {Event} from './event.model';
import {Error} from './error.model';

export enum LogActionTypes {
  AddError = '[Log] Add Error',
  AddEvent = '[Log] Add Event',
}

export class AddEvent implements Action {
  readonly type = LogActionTypes.AddEvent;

  constructor(public payload: { event: Event }) {
  }
}

export class AddError implements Action {
  readonly type = LogActionTypes.AddError;

  constructor(public payload: { error: Error }) {
  }
}

export type LogActions =
  AddError
  | AddEvent;
