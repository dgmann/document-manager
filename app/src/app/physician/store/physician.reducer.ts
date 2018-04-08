import {uniq, without} from 'lodash-es';
import {RequiredAction} from "../../store";
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

    case PhysicianActionTypes.SetRecord:
      const change = {
        escalatedIds: without(state.escalatedIds, action.payload.id),
        reviewIds: without(state.reviewIds, action.payload.id),
        otherIds: without(state.otherIds, action.payload.id),
      };
      if (action.payload.requiredAction !== RequiredAction.NONE) {
        change[action.payload.requiredAction + 'Ids'] = uniq([...state[action.payload.requiredAction + 'Ids'], action.payload.id]);
      }
      return Object.assign({}, state, change);

    default:
      return state;
  }
}
