import {createEntityAdapter, EntityAdapter, EntityState, Update} from '@ngrx/entity';
import {uniq} from 'lodash-es';
import union from 'lodash-es/union';
import without from 'lodash-es/without';
import {RecordActions, RecordActionTypes} from './record.actions';
import {Record} from './record.model';

export interface State extends EntityState<Record> {
  // additional entities state properties
  invalidIds: string[];
  inboxIds: string[];
  escalatedIds: string[];
  reviewIds: string[];
  otherIds: string[];
  doneIds: string[];
  isLoading: boolean;
}

export const adapter: EntityAdapter<Record> = createEntityAdapter<Record>();

export const initialState: State = adapter.getInitialState({
  // additional entity state properties
  invalidIds: [],
  inboxIds: [],
  escalatedIds: [],
  reviewIds: [],
  otherIds: [],
  doneIds: [],
  isLoading: false
});

export function reducer(state = initialState,
                        action: RecordActions): State {
  switch (action.type) {
    case RecordActionTypes.LoadRecords:
      return {
        ...state,
        isLoading: true
      };

    case RecordActionTypes.LoadRecordsFail:
      return {
        ...state,
        isLoading: false
      };

    case RecordActionTypes.LoadRecordsSuccess:
      let stateWithRecord = adapter.addMany(action.payload.records, state);
      stateWithRecord = clearIdsFromState(action.payload.records, stateWithRecord);
      return {
        ...addToStatus(action.payload.records, stateWithRecord),
        isLoading: false
      };

    case RecordActionTypes.UpdateRecordSuccess:
      const updatedState = adapter.updateOne(action.payload.record, state);
      const updatedStateWithoutInvalidId = {
        ...updatedState,
        invalidIds: without(updatedState.invalidIds, action.payload.record.id + '')
      };
      const statusChanges = recordUpdatesToStatusChanges([action.payload.record]);
      const updatedStateWithoutIds = clearIdsFromState(statusChanges, updatedStateWithoutInvalidId);
      return addToStatus(statusChanges, updatedStateWithoutIds);

    case RecordActionTypes.DeleteRecordSuccess:
      const stateWithoutRecord = adapter.removeOne(action.payload.id, state);
      const newState = {
        ...stateWithoutRecord,
        invalidIds: without(stateWithoutRecord.invalidIds, action.payload.id)
      };
      return clearIdsFromState([{id: action.payload.id}], newState);

    case RecordActionTypes.ClearRecords:
      return initialState;

    case RecordActionTypes.UpdateRecord:
      return {
        ...state,
        invalidIds: union(state.invalidIds, [action.payload.record.id as string])
      };

    case RecordActionTypes.UpdatePages:
    case RecordActionTypes.DeleteRecord:
      return {
        ...state,
        invalidIds: union(state.invalidIds, [action.payload.id])
      };

    default:
      return state;
  }
}

function addToStatus(records: StatusChange[], state: State) {
  const change = Object.assign({}, state);
  records.filter(r => !!r.status).forEach(record => change[record.status + 'Ids'] = uniq([...change[record.status + 'Ids'], record.id]));
  return change;
}

function clearIdsFromState(records: StatusChange[], state: State) {
  const ids = records.map(record => record.id);
  const change = {
    inboxIds: without(state.inboxIds, ...ids),
    escalatedIds: without(state.escalatedIds, ...ids),
    reviewIds: without(state.reviewIds, ...ids),
    otherIds: without(state.otherIds, ...ids),
    doneIds: without(state.doneIds, ...ids)
  };
  return {
    ...state,
    ...change
  };
}

interface StatusChange {
  id: string;
  status?: string;
}

function recordUpdatesToStatusChanges(records: Update<Record>[]) {
  return records.filter(r => !!r.changes.status).map(r => ({
    id: r.id,
    status: r.changes.status
  }) as StatusChange);
}

export const {
  selectIds,
  selectEntities,
  selectAll,
  selectTotal,
} = adapter.getSelectors();
