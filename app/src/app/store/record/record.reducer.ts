import { createEntityAdapter, EntityAdapter, EntityState } from '@ngrx/entity';
import union from 'lodash-es/union';
import without from 'lodash-es/without';
import { Record } from './record.model';
import { RecordActions, RecordActionTypes } from './record.actions';

export interface State extends EntityState<Record> {
  // additional entities state properties
  invalidIds: string[];
}

export const adapter: EntityAdapter<Record> = createEntityAdapter<Record>();

export const initialState: State = adapter.getInitialState({
  // additional entity state properties
  invalidIds: []
});

export function reducer(state = initialState,
                        action: RecordActions): State {
  switch (action.type) {
    case RecordActionTypes.LoadRecordsSuccess: {
      return adapter.addMany(action.payload.records, state);
    }

    case RecordActionTypes.UpdateRecordSuccess: {
      let s = adapter.updateOne(action.payload.record, state);
      s.invalidIds = without(s.invalidIds, action.payload.record.id + '');
      return s;
    }

    case RecordActionTypes.DeleteRecordSuccess: {
      let s = adapter.removeOne(action.payload.id, state);
      s.invalidIds = without(s.invalidIds, action.payload.id);
      return s;
    }

    case RecordActionTypes.ClearRecords: {
      return adapter.removeAll(state);
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

export const {
  selectIds,
  selectEntities,
  selectAll,
  selectTotal,
} = adapter.getSelectors();
