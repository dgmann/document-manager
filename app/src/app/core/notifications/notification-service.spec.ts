import {NotificationService} from './notification-service';
import createSpyObj = jasmine.createSpyObj;

describe('NotificationService', () => {
  let service: NotificationService;
  let snackBarService;
  let ngZone;

  beforeEach(() => {
    snackBarService = createSpyObj(['openFromComponent']);
    ngZone = createSpyObj(['runOutsideAngular', 'run']);
    service = new NotificationService(snackBarService, ngZone);
  });

  it('should create', () => {
    expect(service).toBeDefined();
  });
});

