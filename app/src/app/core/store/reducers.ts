import { ActionReducerMap, MetaReducer } from '@ngrx/store';
import { environment } from '../../../environments/environment';
import * as fromRecord from './record/record.reducer';

export interface State {
  records: fromRecord.State;
}

export const reducers: ActionReducerMap<State> = {
  records: fromRecord.reducer,
};


export const metaReducers: MetaReducer<State>[] = !environment.production ? [] : [];
