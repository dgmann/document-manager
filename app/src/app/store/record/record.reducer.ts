import {createEntityAdapter, EntityAdapter, EntityState, Update} from '@ngrx/entity';
import {uniq} from "lodash-es";
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
}

export const adapter: EntityAdapter<Record> = createEntityAdapter<Record>();

export const initialState: State = adapter.getInitialState({
  // additional entity state properties
  invalidIds: [],
  inboxIds: [],
  escalatedIds: [],
  reviewIds: [],
  otherIds: [],
  doneIds: []
});

export function reducer(state = initialState,
                        action: RecordActions): State {
  switch (action.type) {
    case RecordActionTypes.LoadRecordsSuccess: {
      let s = adapter.addMany(action.payload.records, state);
      s = clearIdsFromState(action.payload.records, s);
      return addToStatus(action.payload.records, s);
    }

    case RecordActionTypes.UpdateRecordSuccess: {
      let s = adapter.updateOne(action.payload.record, state);
      s.invalidIds = without(s.invalidIds, action.payload.record.id + '');
      const statusChanges = recordUpdatesToStatusChanges([action.payload.record]);
      s = clearIdsFromState(statusChanges, s);
      return addToStatus(statusChanges, s);
    }

    case RecordActionTypes.DeleteRecordSuccess: {
      let s = adapter.removeOne(action.payload.id, state);
      s.invalidIds = without(s.invalidIds, action.payload.id);
      return clearIdsFromState([{id: action.payload.id}], s);
    }

    case RecordActionTypes.ClearRecords: {
      let s = adapter.removeAll(state);
      return Object.assign({}, s, {
        inboxIds: [],
        escalatedIds: [],
        reviewIds: [],
        otherIds: [],
        doneIds: []
      })
    }

    case RecordActionTypes.UpdateRecord: {
      return Object.assign({}, state, {
        invalidIds: union(state.invalidIds, [action.payload.record.id])
      });
    }

    case RecordActionTypes.UpdatePages:
    case RecordActionTypes.DeleteRecord: {
      return Object.assign({}, state, {
        invalidIds: union(state.invalidIds, [action.payload.id])
      });
    }

    default: {
      return state;
    }
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
  return Object.assign({}, state, change);
}

interface StatusChange {
  id: string,
  status?: string
}

function recordUpdatesToStatusChanges(records: Update<Record>[]) {
  return records.filter(r => !!r.changes.status).map(r => <StatusChange>{
    id: r.id,
    status: r.changes.status
  });
}

export const {
  selectIds,
  selectEntities,
  selectAll,
  selectTotal,
} = adapter.getSelectors();
