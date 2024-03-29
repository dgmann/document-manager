import {PhysicianActions, PhysicianActionTypes} from './physician.actions';

export interface State {
  selectedIds: string[];
  escalatedIds: string[];
  reviewIds: string[];
  otherIds: string[];
  loading: boolean;
  synced: boolean;
}

export const initialState: State = {
  selectedIds: [],
  escalatedIds: [],
  reviewIds: [],
  otherIds: [],
  loading: false,
  synced: false,
};

export function reducer(state = initialState, action: PhysicianActions): State {
  switch (action.type) {

    case PhysicianActionTypes.LoadRecords:
      return Object.assign({}, state, {
        loading: true
      });

    case PhysicianActionTypes.SelectRecords:
      return Object.assign({}, state, {
        selectedIds: action.payload.ids
      });

    default:
      return state;
  }
}
