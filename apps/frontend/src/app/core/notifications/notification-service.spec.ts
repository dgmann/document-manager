import {NotificationService} from './notification-service';

describe('NotificationService', () => {
  let service: NotificationService;
  let snackBarService;
  let ngZone;

  beforeEach(() => {
    snackBarService = {
      openFromComponent: jest.fn()
    };
    ngZone = {
      runOutsideAngular: jest.fn(),
      run: jest.fn()
    };
    service = new NotificationService(snackBarService, ngZone);
  });

  it('should create', () => {
    expect(service).toBeDefined();
  });
});

