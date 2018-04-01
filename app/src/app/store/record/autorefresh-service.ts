import {Injectable} from "@angular/core";
import {Store} from "@ngrx/store";
import {filter} from "rxjs/operators";
import {environment} from "../../../environments/environment"
import {NotificationService} from "../../shared/notification-service";
import {NotificationMessage, NotificationMessageType, WebsocketService} from "../../shared/websocket-service";
import {State} from "../reducers";
import {DeleteRecordSuccess, LoadRecordsSuccess, UpdateRecordSuccess} from "./record.actions";

@Injectable()
export class AutorefreshService {
  constructor(private store: Store<State>,
              private websocketService: WebsocketService,
              private notificationService: NotificationService) {
  }

  start() {
    let ws = this.websocketService.create(environment.websocket);
    const filterRecordEvents = filter((event: NotificationMessage) =>
      event.type == NotificationMessageType.Created
      || event.type == NotificationMessageType.Updated
      || event.type == NotificationMessageType.Deleted);
    ws.pipe(filterRecordEvents).subscribe(message => {
      switch (message.type) {
        case NotificationMessageType.Created:
          this.store.dispatch(new LoadRecordsSuccess({
            records: [message.data]
          }));
          this.notificationService.publish({
            message: "Neues Dokument hinzugefügt",
            timestamp: message.timestamp,
            record: message.data
          });
          break;
        case NotificationMessageType.Updated:
          this.store.dispatch(new UpdateRecordSuccess({
            record: {
              id: message.data.id as string, changes: message.data
            }
          }));
          this.notificationService.publish({
            message: "Änderungen gespeichert",
            timestamp: message.timestamp,
            record: message.data
          });
          break;
        case NotificationMessageType.Deleted:
          this.store.dispatch(new DeleteRecordSuccess({id: message.data.id as string}));
          this.notificationService.publish({
            message: "Dokument gelöscht",
            timestamp: message.timestamp,
            record: message.data
          });
          break;
      }
    });
  }
}