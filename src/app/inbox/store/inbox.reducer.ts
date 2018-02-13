import { InboxActions, InboxActionTypes } from './inbox.actions';

export interface State {
  selectedIds: string[];
  loading: boolean;
  synced: boolean;
}

export const initialState: State = {
  selectedIds: [],
  loading: false,
  synced: false,
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

    default:
      return state;
  }
}
