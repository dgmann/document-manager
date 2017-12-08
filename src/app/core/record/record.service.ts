import { Store } from "@ngrx/store";
import * as fromRecords from "./record.selectors";
import { State } from "./record.reducer";
import { Record } from "./record.model";
import { AddRecords, LoadInbox } from "./record.actions";

export class RecordService {
  private state: Store<State>;

  constructor(private store: Store<State>) {
    this.state = this.store.select(fromRecords.selectRecordState);
  }

  public get(id: string) {
    return this.state.select(fromRecords.selectRecord(id))
  }

  public getAll() {
    return this.state.select(fromRecords.selectAllRecords);
  }

  public getSelected() {
    return this.state.select(fromRecords.selectCurrentRecord);
  }

  public add(records: Record[]) {
    return this.state.dispatch(new AddRecords({ records: records }))
  }

  public load() {
    this.state.dispatch(new LoadInbox())
  }
}
