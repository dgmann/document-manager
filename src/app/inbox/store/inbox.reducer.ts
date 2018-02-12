import {InboxActions, InboxActionTypes} from './inbox.actions';

export interface State {
  ids: string[];
  selectedIds: string[];
}

export const initialState: State = {
  ids: [],
  selectedIds: []
};

export function reducer(state = initialState, action: InboxActions): State {
  switch (action.type) {

    case InboxActionTypes.AddRecords:
      return Object.assign({}, state, {
        ids: [...state.ids, ...action.payload.ids]
      });
    case InboxActionTypes.RemoveRecords:
      return Object.assign({}, state, {
        ids: state.ids.filter(i => action.payload.ids.indexOf(i) < 0)
      });

    case InboxActionTypes.SelectRecords:
      return Object.assign({}, state, {
        selectedIds: action.payload.ids
      });

    default:
      return state;
  }
}
