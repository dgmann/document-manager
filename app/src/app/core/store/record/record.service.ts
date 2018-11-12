import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { select, Store } from "@ngrx/store";
import { map } from "rxjs/operators";
import { environment } from "../../../../environments/environment";
import { State } from "../reducers";
import { DeleteRecord, LoadRecords, UpdatePages, UpdateRecord } from "./record.actions";
import { PageUpdate, Record } from "./record.model";
import { selectAllRecords, selectInvalidIds, selectIsLoading, selectRecordEntities } from "./record.selectors";
import { Observable } from "rxjs";
import { ActionType, GenericEvent, NotificationService, RecordEvent } from "../../notification-service";

@Injectable({
  providedIn: "root"
})
export class RecordService {
  public records: Observable<Record[]>;
  public invalidIds: Observable<string[]>;
  public isLoading$: Observable<boolean>;

  constructor(private store: Store<State>, private http: HttpClient, private notifications: NotificationService) {
    this.records = this.store.pipe(select(selectAllRecords));
    this.invalidIds = this.store.pipe(select(selectInvalidIds));
    this.isLoading$ = this.store.pipe(select(selectIsLoading));
  }

  public load(query: { [param: string]: string }) {
    this.store.dispatch(new LoadRecords({query: query}))
  }

  public find(id: string) {
    return this.store.pipe(select(selectRecordEntities), map(entities => entities[id]));
  }

  public delete(id: string) {
    this.store.dispatch(new DeleteRecord({id: id}))
  }

  public update(id: string, changes: any) {
    this.store.dispatch(new UpdateRecord({record: {id: id, changes: changes}}))
  }

  public updatePages(id: string, pages: PageUpdate[]) {
    this.store.dispatch(new UpdatePages({id: id, updates: pages}))
  }

  public upload(pdf) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: "PDF wird hochgeladen..."
    }));

    const formData = new FormData();
    formData.append('pdf', pdf);
    formData.append('sender', "Client");
    this.http.post<Record>(environment.api + "/records", formData)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: "PDF hochgeladen",
        record: record
    })));
  }

  public append(sourceId: string, targetId: string) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: "PDF wird angehängt..."
    }));
    this.http.post<Record>(`${environment.api}/records/${sourceId}/append/${targetId}`, null)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: "PDF angehängt",
        record: record
      })));
  }

  public reset(id: string) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: "Befund wird zurückgesetzt..."
    }));
    this.http.put<Record>(`${environment.api}/records/${id}/reset`, null)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: "Befund zurückgesetzt",
        record: record
    })));
  }
}
