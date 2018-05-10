import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {selectInboxRecords} from "../store";
import {selectMultiselect, selectSelectedIds, selectSelectedRecords, selectUnreadRecords, State} from "./reducers";
import {SelectRecords, SetMultiSelect} from "./store/inbox.actions";

@Injectable()
export class InboxService {

  constructor(private store: Store<State>) {
  }

  public all() {
    return this.store.pipe(select(selectInboxRecords));
  }

  public getSelectedIds() {
    return this.store.pipe(select(selectSelectedIds))
  }

  public getSelectedRecords() {
    return this.store.pipe(select(selectSelectedRecords));
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids: ids}))
  }

  public getUnreadRecords() {
    return this.store.pipe(select(selectUnreadRecords));
  }

  public setMultiselect(value: boolean) {
    this.store.dispatch(new SetMultiSelect({multiselect: value}));
  }

  public getMultiselect() {
    return this.store.pipe(select(selectMultiselect))
  }
}
