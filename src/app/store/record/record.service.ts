import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Dictionary } from "@ngrx/entity/src/models";
import { Store } from "@ngrx/store";
import { map } from "rxjs/operators";
import { environment } from "../../../environments/environment";
import { NotificationService } from "../../shared/notification-service";
import { State } from "../reducers";
import { DeleteRecord, LoadRecords, UpdateRecord } from "./record.actions";
import { Record } from "./record.model";
import { selectAllRecords, selectRecordEntities } from "./record.selectors";

@Injectable()
export class RecordService {
  constructor(private store: Store<State>, private http: HttpClient, private notificationService: NotificationService) {
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

  public delete(id: string) {
    this.store.dispatch(new DeleteRecord({id: id}))
  }

  public update(id: string, changes: any) {
    this.store.dispatch(new UpdateRecord({record: {id: id, changes: changes}}))
  }

  public upload(pdf) {
    const formData = new FormData();
    formData.append('pdf', pdf);
    formData.append('sender', "Client");
    this.http.post(environment.api + "/records", formData).subscribe(() => this.notificationService.publish({
      timestamp: new Date(),
      message: "PDF hochgeladen"
    }))
  }

  public append(sourceId: string, targetId: string) {
    this.http.post(`${environment.api}/records/${sourceId}/append/${targetId}`, null).subscribe(() => this.notificationService.publish({
      timestamp: new Date(),
      message: "PDF angehängt"
    }))
  }
}
