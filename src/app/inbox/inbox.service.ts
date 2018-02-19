import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {filter, map} from "rxjs/operators";
import {Record, RecordService} from "../store";
import {selectSelectedIds, State} from "./reducers";
import {SelectRecords} from "./store/inbox.actions";

@Injectable()
export class InboxService {

  inboxFilter: any;

  constructor(private store: Store<State>, private recordService: RecordService) {
    this.inboxFilter = (record: Record) => (!record.date || !record.patientId || !record.tags) && !record.requiredAction
  }

  public load() {
    this.recordService.load({inbox: "true"})
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
