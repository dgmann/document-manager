import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {filter, map} from "rxjs/operators";
import {RecordService} from "../store";
import {selectSelectedIds, State} from "./reducers";
import {SelectRecords} from "./store/inbox.actions";

@Injectable()
export class InboxService {

  inboxFilter: any;

  constructor(private store: Store<State>, private recordService: RecordService) {
    this.inboxFilter = record => !record.escalated && !record.processed
  }

  public load() {
    this.recordService.load({escalated: "false", processed: "false"})
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

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids: ids}))
  }
}
