import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {filter, map} from "rxjs/operators";
import {Record, RecordService, Status} from "../store";
import {selectMultiselect, selectSelectedIds, selectSelectedRecords, selectUnreadRecords, State} from "./reducers";
import {SelectRecords, SetMultiSelect} from "./store/inbox.actions";

@Injectable()
export class InboxService {

  inboxFilter: any;

  constructor(private store: Store<State>, private recordService: RecordService) {
    this.inboxFilter = (record: Record) => !record.status || record.status == Status.INBOX;
  }

  public all() {
    return this.recordService.all().pipe(map(records => records.filter(this.inboxFilter)))
  }

  public find(id: string) {
    return this.recordService.find(id).pipe(filter(this.inboxFilter))
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
