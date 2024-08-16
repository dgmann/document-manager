import {Injectable} from '@angular/core';
import {select, Store} from '@ngrx/store';
import {
  Record,
  RecordService,
  selectEscalatedRecords,
  selectOtherRecords,
  selectReviewRecords,
  Status
} from '../core/records';
import {selectSelectedIds, selectSelectedRecords, State} from './reducers';
import {SelectRecords} from './store/physician.actions';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class PhysicianService {
  public selectedRecords$: Observable<Record[]>;
  public selectedIds$: Observable<string[]>;
  public escalatedRecords$: Observable<Record[]>;
  public reviewRecords$: Observable<Record[]>;
  public otherRecords$: Observable<Record[]>;

  constructor(private store: Store<State>, private recordService: RecordService) {
    this.selectedRecords$ = this.store.pipe(select(selectSelectedRecords));
    this.selectedIds$ = this.store.pipe(select(selectSelectedIds));
    this.escalatedRecords$ = this.store.pipe(select(selectEscalatedRecords));
    this.reviewRecords$ = this.store.pipe(select(selectReviewRecords));
    this.otherRecords$ = this.store.pipe(select(selectOtherRecords));
  }

  public load() {
    this.recordService.load({status: Status.REVIEW});
    this.recordService.load({status: Status.ESCALATED});
    this.recordService.load({status: Status.OTHER});
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids}));
  }
}
