import { NotificationEvent, NotificationService } from "./notification-service";
import { async } from "@angular/core/testing";
import createSpyObj = jasmine.createSpyObj;

describe("NotificationService", () => {
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

  it('should publish events', async(() => {
    const event = createEvent("Test");
    service.publish(event);
    service.events$.subscribe(e => expect(e).toEqual(event));
  }));
});

const createEvent = (message: string) => ({
  payload: {
    message,
    timestamp: new Date()
  }
} as NotificationEvent);
