import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {RecordService} from "../store";
import {
  selectEscalatedRecords,
  selectOtherRecords,
  selectReviewRecords,
  selectSelectedIds,
  selectSelectedRecords,
  State
} from "./reducers";
import {SelectRecords} from "./store/physician.actions";

@Injectable()
export class PhysicianService {

  constructor(private store: Store<State>, private recordService: RecordService) {
  }

  public getSelectedRecords() {
    return this.store.pipe(select(selectSelectedRecords));
  }

  public getSelectedIds() {
    return this.store.pipe(select(selectSelectedIds));
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids: ids}))
  }

  public getEscalated() {
    return this.store.pipe(select(selectEscalatedRecords));
  }

  public getToReview() {
    return this.store.pipe(select(selectReviewRecords));
  }

  public getOther() {
    return this.store.pipe(select(selectOtherRecords));
  }
}
