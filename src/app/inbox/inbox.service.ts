import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {selectIds, selectSelectedIds, State} from "./reducers";
import {AddRecords, RemoveRecords, SelectRecords} from "./store/inbox.actions";

@Injectable()
export class InboxService {
  constructor(private store: Store<State>) {
  }

  public getSelectedIds() {
    return this.store.pipe(select(selectSelectedIds))
  }

  public selectIds(ids: string[]) {
    this.store.dispatch(new SelectRecords({ids: ids}))
  }

  public getIds() {
    return this.store.pipe(select(selectIds))
  }

  public addIds(ids: string[]) {
    this.store.dispatch(new AddRecords({ids: ids}))
  }

  public removeIds(ids: string[]) {
    this.store.dispatch(new RemoveRecords({ids: ids}))
  }
}
