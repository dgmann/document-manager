import {Injectable} from "@angular/core";
import {MatSnackBar} from "@angular/material";
import {Subject} from "rxjs/Subject";
import {Record} from "../store";

@Injectable()
export class NotificationService {
  private events = new Subject<NotificationEvent>();

  constructor(public snackbar: MatSnackBar) {
  }

  publish(event: NotificationEvent) {
    this.events.next(event);
  }

  logToConsole() {
    this.events.subscribe(event => console.log(`${event.timestamp} ${event.message}: ${event.record && event.record.id}`))
  }

  logToSnackBar() {
    this.events.subscribe(event => this.snackbar.open(event.message, '', {
      duration: 3000
    }));
  }
}

export interface NotificationEvent {
  message: string;
  timestamp: Date;
  record?: Record;
}
