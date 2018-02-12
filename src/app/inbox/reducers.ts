import {ActionReducerMap, createSelector, MetaReducer} from '@ngrx/store';
import {environment} from '../../environments/environment';
import * as fromInbox from './store/inbox.reducer';

export interface State {
  inbox: fromInbox.State;
}

export const reducers: ActionReducerMap<State> = {
  inbox: fromInbox.reducer,
};


export const metaReducers: MetaReducer<State>[] = !environment.production ? [] : [];

export const selectFeature = (state: State) => state.inbox;
export const selectIds = createSelector(selectFeature, (state: fromInbox.State) => state.ids);
export const selectSelectedIds = createSelector(selectFeature, (state: fromInbox.State) => state.selectedIds);
