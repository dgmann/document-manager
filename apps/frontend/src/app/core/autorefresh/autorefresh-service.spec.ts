import {AutorefreshService} from './autorefresh-service';
import {NotificationMessage, NotificationMessageType, NotificationTopic} from '../notifications/websocket-service';
import {ActionType, RecordEvent} from '../notifications/notification-service';
import {of, throwError} from 'rxjs';

describe('AutoRefreshService', () => {
  let service: AutorefreshService;
  let notificationService;
  let websocketService;
  let store;
  let configService;

  beforeEach(() => {
    notificationService = {
      publish: jest.fn()
    };
    websocketService = {
      create: jest.fn().mockReturnValue(of(null))
    };
    store = {
      pipe: jest.fn(),
      dispatch: jest.fn()
    };
    configService = {
      getNotificationWebsocketUrl: jest.fn().mockReturnValue('http://test.com')
    };
    service = new AutorefreshService(store, websocketService, notificationService, configService);
  });

  it('should create service', () => {
    expect(service.webSocket$).toBeDefined();
  });

  it('should handle create messages', () => {
    const createMessage: NotificationMessage = {
      id: '1',
      timestamp: new Date(),
      type: NotificationMessageType.Created,
      topic: NotificationTopic.Records
    };
    service.webSocket$ = of(createMessage);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.ADDED,
      message: 'Neues Dokument hinzugefügt',
      timestamp: createMessage.timestamp,
      id: createMessage.id
    });

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalled();
  });

  it('should handle update messages', () => {
    const message: NotificationMessage = {
      id: '1',
      timestamp: new Date(),
      type: NotificationMessageType.Updated,
      topic: NotificationTopic.Records
    };
    service.webSocket$ = of(message);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.UPDATED,
      message: 'Änderungen gespeichert',
      timestamp: message.timestamp,
      id: message.id
    });

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalled();
  });

  it('should handle delete messages', () => {
    const message: NotificationMessage = {
      id: '1',
      timestamp: new Date(),
      type: NotificationMessageType.Deleted,
      topic: NotificationTopic.Records
    };
    service.webSocket$ = of(message);

    service.start();

    const expectedNotification = new RecordEvent({
      type: ActionType.DELETED,
      message: 'Dokument gelöscht',
      timestamp: message.timestamp,
      id: message.id
    });

    expect(notificationService.publish).toHaveBeenCalledWith(expectedNotification);
    expect(store.dispatch).toHaveBeenCalled();
  });

  it('should handle error', () => {
    service.webSocket$ = throwError(null);
    service.start();

    expect(notificationService.publish).toHaveBeenCalled();
  });
});
