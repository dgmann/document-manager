import {Injectable} from '@angular/core';
import {select, Store} from '@ngrx/store';
import {Record, RecordService, selectInboxIds, selectInboxRecords, Status} from '../core/records';
import {selectSelectedIds, selectSelectedRecords, State} from './reducers';
import {SelectRecords} from './store/inbox.actions';
import {Observable} from 'rxjs';
import {debounceTime, take} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class InboxService {
  public allInboxRecords$: Observable<Record[]>;
  public allInboxRecordIds$: Observable<string[]>;
  public selectedIds$: Observable<string[]>;
  public selectedRecords$: Observable<Record[]>;
  public isLoading$: Observable<boolean>;

  constructor(private store: Store<State>,
              private recordService: RecordService) {
    this.allInboxRecords$ = this.store.pipe(select(selectInboxRecords));
    this.allInboxRecordIds$ = this.store.pipe(select(selectInboxIds));
    this.selectedIds$ = this.store.pipe(select(selectSelectedIds));
    this.selectedRecords$ = this.store.pipe(select(selectSelectedRecords));
    this.isLoading$ = this.recordService.isLoading$.pipe(debounceTime(1000));
  }

  public loadRecords() {
    this.recordService.load({status: Status.INBOX});
  }

  public upload(pdf) {
    this.recordService.upload(pdf);
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids}));
  }

  public deleteSelectedRecords() {
    const deleteMethod = (id: string) => this.recordService.delete(id);
    this.doForAllSelectedRecords(deleteMethod);
  }

  public updateSelectedRecords(changes: any) {
    const update = (id: string) => this.recordService.update(id, changes);
    this.doForAllSelectedRecords(update);
  }

  private doForAllSelectedRecords(callback: (id: string) => any) {
    this.selectedIds$.pipe(take(1)).subscribe(ids => ids.forEach(id => callback(id)));
  }
}
