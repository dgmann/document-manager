import { Injectable } from "@angular/core";
import { Store } from "@ngrx/store";
import { State } from "../reducers";
import { selectAllRecords, selectRecordEntities } from "./record.selectors";
import { LoadRecords, UpdateRecord } from "./record.actions";
import { Record } from "./record.model";
import { Dictionary } from "@ngrx/entity/src/models";
import { map } from "rxjs/operators";

@Injectable()
export class RecordService {
  constructor(private store: Store<State>) {
  }

  public load() {
    this.store.dispatch(new LoadRecords(null))
  }

  public find(id: string) {
    return this.store.select<Dictionary<Record>>(selectRecordEntities).pipe<Record>(map(entities => entities[id]));
  }

  public all() {
    return this.store.select<Record[]>(selectAllRecords)
  }

  public update(id: string, changes: any) {
    this.store.dispatch(new UpdateRecord({record: {id: id, changes: changes}}))
  }
}
