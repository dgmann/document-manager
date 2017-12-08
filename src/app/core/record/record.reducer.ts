import { createEntityAdapter, EntityAdapter, EntityState } from "@ngrx/entity";
import { Record } from "./record.model";
import * as RecordActions from "./record.actions"
import { PayloadAction } from "../payload-action";

export interface State extends EntityState<Record> {
  selectedRecordId: string | null;
}

export function sortByDate(a: Record, b: Record): number {
  return getTime(a.date) - getTime(b.date);
}

function getTime(date?: Date) {
  return date != null ? date.getTime() : 0;
}

export const adapter: EntityAdapter<Record> = createEntityAdapter<Record>({
  sortComparer: sortByDate,
});


export const initialState: State = adapter.getInitialState({
  selectedRecordId: null
});

export function reducer(state = initialState, action: PayloadAction
): State {
  switch (action.type) {
    case RecordActions.ADD_RECORD: {
      return adapter.addOne(action.payload.record, state);
    }

    case RecordActions.ADD_RECORDS: {
      return adapter.addMany(action.payload.records, state);
    }

    case RecordActions.UPDATE_RECORD: {
      return adapter.updateOne(action.payload.record, state);
    }

    case RecordActions.UPDATE_RECORDS: {
      return adapter.updateMany(action.payload.records, state);
    }

    case RecordActions.DELETE_RECORD: {
      return adapter.removeOne(action.payload.id, state);
    }

    case RecordActions.DELETE_RECORDS: {
      return adapter.removeMany(action.payload.ids, state);
    }

    case RecordActions.LOAD_RECORDS: {
      return adapter.addAll(action.payload.records, state);
    }

    case RecordActions.CLEAR_RECORDS: {
      return adapter.removeAll({ ...state, selectedRecordId: null });
    }

    default: {
      return state;
    }
  }
}

export const getSelectedRecordId = (state: State) => state.selectedRecordId;
