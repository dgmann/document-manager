import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {RecordService, RequiredAction} from "../store";
import {
  selectEscalatedRecords,
  selectOtherRecords,
  selectReviewRecords,
  selectSelectedRecords,
  State
} from "./reducers";
import {SelectRecords} from "./store/physician.actions";

@Injectable()
export class PhysicianService {

  constructor(private store: Store<State>, private recordService: RecordService) {
  }

  public load(action: RequiredAction) {
    this.recordService.load({requiredAction: action})
  }

  public getSelectedRecords() {
    return this.store.pipe(select(selectSelectedRecords));
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
