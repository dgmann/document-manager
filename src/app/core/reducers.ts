import { ActionReducerMap } from '@ngrx/store';
import * as fromRecord from './record/record.reducer';

export interface State {
  records: fromRecord.State;
}

export const reducers: ActionReducerMap<State> = {
  records: fromRecord.reducer
};
