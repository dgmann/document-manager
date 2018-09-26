import { Injectable } from "@angular/core";
import { select, Store } from "@ngrx/store";
import { Record, RecordService, selectInboxIds, selectInboxRecords, Status } from "../core/store";
import { selectMultiselect, selectSelectedIds, selectSelectedRecords, State } from "./reducers";
import { SelectRecords, SetMultiSelect } from "./store/inbox.actions";
import { Observable } from "rxjs";
import { take } from "rxjs/operators";

@Injectable()
export class InboxService {
  public allInboxRecords$: Observable<Record[]>;
  public allInboxRecordIds$: Observable<string[]>;
  public selectedIds$: Observable<string[]>;
  public selectedRecords$: Observable<Record[]>;
  public isMultiSelect$: Observable<boolean>;

  constructor(private store: Store<State>,
              private recordService: RecordService) {
    this.allInboxRecords$ = this.store.pipe(select(selectInboxRecords));
    this.allInboxRecordIds$ = this.store.pipe(select(selectInboxIds));
    this.selectedIds$ = this.store.pipe(select(selectSelectedIds));
    this.selectedRecords$ = this.store.pipe(select(selectSelectedRecords));
    this.isMultiSelect$ = this.store.pipe(select(selectMultiselect));
  }

  public loadRecords() {
    this.recordService.load({status: Status.INBOX});
  }

  public upload(pdf) {
    this.recordService.upload(pdf);
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids: ids}))
  }

  public setMultiselect(value: boolean) {
    this.store.dispatch(new SetMultiSelect({multiselect: value}));
  }

  public deleteSelectedRecords() {
    this.doForAllSelectedRecords(this.recordService.delete);
  }

  public updateSelectedRecords(changes: any) {
    let update = (id: string) => this.recordService.update(id, changes);
    this.doForAllSelectedRecords(update);
  }

  private doForAllSelectedRecords(callback: (id: string) => void) {
    this.selectedIds$.pipe(take(1)).subscribe(ids => ids.forEach(id => callback(id)));
  }
}
