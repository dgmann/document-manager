import { createEntityAdapter, EntityAdapter, EntityState } from '@ngrx/entity';
import { Record } from './record.model';
import { RecordActions, RecordActionTypes } from './record.actions';

export interface State extends EntityState<Record> {
  // additional entities state properties
}

export const adapter: EntityAdapter<Record> = createEntityAdapter<Record>();

export const initialState: State = adapter.getInitialState({
  // additional entity state properties
});

export function reducer(state = initialState,
                        action: RecordActions): State {
  switch (action.type) {
    case RecordActionTypes.AddRecord: {
      return adapter.addOne(action.payload.record, state);
    }

    case RecordActionTypes.AddRecords: {
      return adapter.addMany(action.payload.records, state);
    }

    case RecordActionTypes.UpdateRecord: {
      return adapter.updateOne(action.payload.record, state);
    }

    case RecordActionTypes.UpdateRecords: {
      return adapter.updateMany(action.payload.records, state);
    }

    case RecordActionTypes.DeleteRecord: {
      return adapter.removeOne(action.payload.id, state);
    }

    case RecordActionTypes.DeleteRecords: {
      return adapter.removeMany(action.payload.ids, state);
    }

    case RecordActionTypes.ClearRecords: {
      return adapter.removeAll(state);
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
