import union from 'lodash-es/union';
import without from 'lodash-es/without';
import {InboxActions, InboxActionTypes} from './inbox.actions';

export interface State {
  selectedIds: string[];
  loading: boolean;
  synced: boolean;
  unreadIds: string[];
  multiselect: boolean;
}

export const initialState: State = {
  selectedIds: [],
  loading: false,
  synced: false,
  unreadIds: [],
  multiselect: false
};

export function reducer(state = initialState, action: InboxActions): State {
  switch (action.type) {

    case InboxActionTypes.LoadRecords:
      return Object.assign({}, state, {
        loading: true
      });

    case InboxActionTypes.SelectRecords:
      return Object.assign({}, state, {
        selectedIds: action.payload.ids
      });
    case InboxActionTypes.AddUnreadRecords:
      return Object.assign({}, state, {
        unreadIds: union(state.unreadIds, action.payload.ids)
      });
    case InboxActionTypes.RemoveUnreadRecords:
      return Object.assign({}, state, {
        unreadIds: without(state.unreadIds, ...action.payload.ids)
      });
    case InboxActionTypes.SetMultiSelect:
      return Object.assign({}, state, {
        multiselect: action.payload.multiselect
      });

    default:
      return state;
  }
}
