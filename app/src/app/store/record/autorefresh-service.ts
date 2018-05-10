import {Injectable} from "@angular/core";
import {Store} from "@ngrx/store";
import {filter, retry} from "rxjs/operators";
import {environment} from "../../../environments/environment"
import {ActionType, GenericEvent, NotificationService, RecordEvent} from "../../shared/notification-service";
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
    ws.pipe(filterRecordEvents, retry()).subscribe(message => {
      switch (message.type) {
        case NotificationMessageType.Created:
          this.store.dispatch(new LoadRecordsSuccess({
            records: [message.data]
          }));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.ADDED,
            message: "Neues Dokument hinzugefügt",
            timestamp: message.timestamp,
            record: message.data
          }));
          break;
        case NotificationMessageType.Updated:
          this.store.dispatch(new UpdateRecordSuccess({
            record: {
              id: message.data.id as string, changes: message.data
            }
          }));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.UPDATED,
            message: "Änderungen gespeichert",
            timestamp: message.timestamp,
            record: message.data
          }));
          break;
        case NotificationMessageType.Deleted:
          this.store.dispatch(new DeleteRecordSuccess({id: message.data.id as string}));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.DELETED,
            message: "Dokument gelöscht",
            timestamp: message.timestamp,
            record: message.data
          }));
          break;
      }
    }, () => {
      this.notificationService.publish(new GenericEvent({
        message: "Verbindung zum Server verloren",
        timestamp: new Date()
      }))
    });
  }
}
