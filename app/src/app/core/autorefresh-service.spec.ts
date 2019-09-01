import {AutorefreshService} from './autorefresh-service';
import {NotificationMessage, NotificationMessageType} from './websocket-service';
import {ActionType, RecordEvent} from './notification-service';
import {DeleteRecordSuccess, LoadRecordsSuccess, UpdateRecordSuccess} from './store';
import {of, throwError} from 'rxjs';
import createSpyObj = jasmine.createSpyObj;

describe('AutoRefreshService', () => {
  let service: AutorefreshService;
  let notificationService;
  let websocketService;
  let store;

  beforeEach(() => {
    notificationService = createSpyObj(['publish']);
    websocketService = createSpyObj(['create']);
    websocketService.create.and.returnValue(of(null));
    store = createSpyObj(['pipe', 'dispatch']);
    service = new AutorefreshService(store, websocketService, notificationService);
  });

  it('should create service', () => {
    expect(service.webSocket$).toBeDefined();
  });

  it('should handle create messages', () => {
    const createMessage: NotificationMessage = {
      data: {},
      timestamp: new Date(),
      type: NotificationMessageType.Created
    };
    service.webSocket$ = of(createMessage);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.ADDED,
      message: 'Neues Dokument hinzugefügt',
      timestamp: createMessage.timestamp,
      record: createMessage.data
    });
    const expectedStoreAction = new LoadRecordsSuccess({
      records: [createMessage.data]
    });

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalledWith(expectedStoreAction);
  });

  it('should handle update messages', () => {
    const message: NotificationMessage = {
      data: {id: '1'},
      timestamp: new Date(),
      type: NotificationMessageType.Updated
    };
    service.webSocket$ = of(message);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.UPDATED,
      message: 'Änderungen gespeichert',
      timestamp: message.timestamp,
      record: message.data
    });
    const expectedStoreAction = new UpdateRecordSuccess({
      record: {
        id: message.data.id as string, changes: message.data
      }
    });

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalledWith(expectedStoreAction);
  });

  it('should handle delete messages', () => {
    const message: NotificationMessage = {
      data: {id: '1'},
      timestamp: new Date(),
      type: NotificationMessageType.Deleted
    };
    service.webSocket$ = of(message);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.DELETED,
      message: 'Dokument gelöscht',
      timestamp: message.timestamp,
      record: message.data
    });
    const expectedStoreAction = new DeleteRecordSuccess({id: message.data.id as string});

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalledWith(expectedStoreAction);
  });

  it('should handle error', () => {
    service.webSocket$ = throwError(null);
    service.start();

    expect(notificationService.publish).toHaveBeenCalled();
  });
});
