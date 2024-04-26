import {LogActions, LogActionTypes} from './log.actions';
import {Event} from './event.model';

export interface State {
  events: Event[];
  errors: any[];
}

export const initialState: State = {
  events: [],
  errors: []
};

export function reducer(state = initialState, action: LogActions): State {
  switch (action.type) {

    case LogActionTypes.AddEvent:
      return {
        ...state,
        events: [action.payload.event, ...state.events]
      };

    case LogActionTypes.AddError:
      return {
        ...state,
        events: [action.payload.error, ...state.errors]
      };

    default:
      return state;
  }
}
