import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {Dictionary} from "@ngrx/entity/src/models";
import {Store} from "@ngrx/store";
import {map} from "rxjs/operators";
import {environment} from "../../../environments/environment";
import {ActionType, NotificationService, RecordEvent} from "../../shared/notification-service";
import {State} from "../reducers";
import {DeleteRecord, LoadRecords, UpdatePages, UpdateRecord} from "./record.actions";
import {PageUpdate, Record} from "./record.model";
import {selectAllRecords, selectInvalidIds, selectRecordEntities} from "./record.selectors";

@Injectable()
export class RecordService {
  constructor(private store: Store<State>, private http: HttpClient, private notificationService: NotificationService) {
  }

  public load(query: { [param: string]: string }) {
    this.store.dispatch(new LoadRecords({query: query}))
  }

  public find(id: string) {
    return this.store.select<Dictionary<Record>>(selectRecordEntities).pipe<Record>(map(entities => entities[id]));
  }

  public all() {
    return this.store.select<Record[]>(selectAllRecords)
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

  public getInvalidIds() {
    return this.store.select<string[]>(selectInvalidIds);
  }

  public upload(pdf) {
    const formData = new FormData();
    formData.append('pdf', pdf);
    formData.append('sender', "Client");
    this.http.post<Record>(environment.api + "/records", formData).subscribe(record => this.notificationService.publish(new RecordEvent({
      type: ActionType.NONE,
      timestamp: new Date(),
      message: "PDF hochladen...",
      record: record
    })));
  }

  public append(sourceId: string, targetId: string) {
    this.http.post<Record>(`${environment.api}/records/${sourceId}/append/${targetId}`, null).subscribe(record => this.notificationService.publish(new RecordEvent({
      type: ActionType.NONE,
      timestamp: new Date(),
      message: "PDF anhängen...",
      record: record
    })))
  }
}
