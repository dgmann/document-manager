import {ActionReducerMap, MetaReducer} from '@ngrx/store';
import {environment} from '../../environments/environment';
import * as fromInbox from './store/inbox.reducer';

export interface State {
  inbox: fromInbox.State;
}

export const reducers: ActionReducerMap<State> = {
  inbox: fromInbox.reducer,
};


export const metaReducers: MetaReducer<State>[] = !environment.production ? [] : [];
