import {Injectable} from '@angular/core';
import {Store} from '@ngrx/store';
import {filter, retry} from 'rxjs/operators';
import {environment} from '@env/environment';
import {DeleteRecordSuccess, LoadRecordsSuccess, Record, State, UpdateRecordSuccess} from './store';
import {ActionType, GenericEvent, NotificationService, RecordEvent} from './notification-service';
import {NotificationMessage, NotificationMessageType, WebsocketService} from './websocket-service';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AutorefreshService {
  public webSocket$: Observable<NotificationMessage>;

  constructor(private store: Store<State>,
              private websocketService: WebsocketService,
              private notificationService: NotificationService) {
    const ws = this.websocketService.create(environment.websocket);
    const filterRecordEvents = filter((event: NotificationMessage) =>
      event.type === NotificationMessageType.Created
      || event.type === NotificationMessageType.Updated
      || event.type === NotificationMessageType.Deleted);
    this.webSocket$ = ws.pipe(filterRecordEvents, retry());
  }

  start() {
    this.webSocket$.subscribe(message => {
      switch (message.type) {
        case NotificationMessageType.Created:
          this.store.dispatch(new LoadRecordsSuccess({
            records: [this.relativeUrlToAbsolute(message.data)]
          }));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.ADDED,
            message: 'Neues Dokument hinzugefügt',
            timestamp: message.timestamp,
            record: message.data
          }));
          break;
        case NotificationMessageType.Updated:
          this.store.dispatch(new UpdateRecordSuccess({
            record: {
              id: message.data.id as string, changes: this.relativeUrlToAbsolute(message.data)
            }
          }));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.UPDATED,
            message: 'Änderungen gespeichert',
            timestamp: message.timestamp,
            record: message.data
          }));
          break;
        case NotificationMessageType.Deleted:
          this.store.dispatch(new DeleteRecordSuccess({id: message.data.id as string}));
          this.notificationService.publish(new RecordEvent({
            type: ActionType.DELETED,
            message: 'Dokument gelöscht',
            timestamp: message.timestamp,
            record: message.data
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

  relativeUrlToAbsolute(record: Record) {
    return {
      ...record,
      pages: record.pages.map(page => ({
        ...page,
        url: `${environment.api}${page.url}`
      }))
    } as Record;
  }
}
