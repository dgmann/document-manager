import { ActionReducerMap, MetaReducer } from '@ngrx/store';
import { environment } from '../../../environments/environment';
import * as fromRecord from './record/record.reducer';
import * as fromLog from './log/log.reducer';


export interface State {
  records: fromRecord.State;
  log: fromLog.State
}

export const reducers: ActionReducerMap<State> = {
  records: fromRecord.reducer,
  log: fromLog.reducer
};


export const metaReducers: MetaReducer<State>[] = !environment.production ? [] : [];
