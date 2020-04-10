import {Injectable} from '@angular/core';
import {Store} from '@ngrx/store';
import {filter, retry} from 'rxjs/operators';
import {LoadRecords} from '@app/core/records';
import {ActionType, GenericEvent, NotificationService, RecordEvent} from '../notifications/notification-service';
import {NotificationMessage, NotificationMessageType, WebsocketService} from '../notifications/websocket-service';
import {Observable} from 'rxjs';
import {ConfigService} from '@app/core/config';
import {State} from '@app/core/store';

@Injectable({
  providedIn: 'root'
})
export class AutorefreshService {
  public webSocket$: Observable<NotificationMessage>;

  constructor(private store: Store<State>,
              private websocketService: WebsocketService,
              private notificationService: NotificationService,
              private config: ConfigService) {
    const ws = this.websocketService.create(this.config.getNotificationWebsocketUrl());
    const filterRecordEvents = filter((event: NotificationMessage) =>
      event.type === NotificationMessageType.Created
      || event.type === NotificationMessageType.Updated
      || event.type === NotificationMessageType.Deleted);
    this.webSocket$ = ws.pipe(filterRecordEvents, retry());
  }

  start() {
    this.webSocket$.subscribe(message => {
      this.store.dispatch(new LoadRecords({query: {id: message.id}}));
      switch (message.type) {
        case NotificationMessageType.Created:
          this.notificationService.publish(new RecordEvent({
            type: ActionType.ADDED,
            message: 'Neues Dokument hinzugefügt',
            timestamp: message.timestamp,
            id: message.id
          }));
          break;
        case NotificationMessageType.Updated:
          this.notificationService.publish(new RecordEvent({
            type: ActionType.UPDATED,
            message: 'Änderungen gespeichert',
            timestamp: message.timestamp,
            id: message.id
          }));
          break;
        case NotificationMessageType.Deleted:
          this.notificationService.publish(new RecordEvent({
            type: ActionType.DELETED,
            message: 'Dokument gelöscht',
            timestamp: message.timestamp,
            id: message.id
          }));
          break;
      }
    }, () => {
      this.notificationService.publish(new GenericEvent({
        message: 'Verbindung zum Server verloren',
        timestamp: new Date()
      }));
    });
  }
}
