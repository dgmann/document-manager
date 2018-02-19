import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {map} from "rxjs/operators";
import {RecordService, RequiredAction} from "../store";
import {selectSelectedIds, selectSelectedRecords, State} from "./reducers";
import {SelectRecords} from "./store/physician.actions";

@Injectable()
export class PhysicianService {

  constructor(private store: Store<State>, private recordService: RecordService) {
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

  public getEscalated() {
    return this.get(record => record.statusType === RequiredAction.ESCALATED)
  }

  public getToReview() {
    return this.get(record => record.statusType === RequiredAction.REVIEW)
  }

  public getOther() {
    return this.get(record => record.statusType === RequiredAction.OTHER)
  }

  private get(filter) {
    return this.recordService.all().pipe(map(records => records.filter(filter)))
  }
}
