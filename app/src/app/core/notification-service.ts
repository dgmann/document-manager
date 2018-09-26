import { Injectable, NgZone } from "@angular/core";
import { MatSnackBar } from "@angular/material";
import { groupBy, map } from "lodash-es";
import { Subject } from "rxjs";
import { bufferTime, filter } from "rxjs/operators";
import { Record } from "./store";
import { EventSnackbarComponent } from "../shared";

@Injectable({
  providedIn: "root"
})
export class NotificationService {
  private _events$ = new Subject<NotificationEvent>();

  public get events$() {
    return this._events$.asObservable();
  }

  constructor(public snackbar: MatSnackBar, private ngZone: NgZone) {
  }

  publish(event: NotificationEvent) {
    this._events$.next(event);
  }

  logToConsole() {
    this.events$.subscribe(event => console.log(event.toString()))
  }

  logToSnackBar() {
    this.ngZone.runOutsideAngular(() => {
      this.events$
        .pipe(bufferTime(500), filter(events => events && events.length > 0))
        .subscribe(events => {
          const grouped = groupBy(events, 'payload.type');
          const messages = map(grouped, (group, key) => {
            const word = group.length == 1 ? "Befund" : "Befunde";
            if (key == ActionType.UPDATED) {
              return new GenericEvent({timestamp: new Date(), message: `${group.length} ${word} geändert.`});
            } else if (key == ActionType.ADDED) {
              return new GenericEvent({timestamp: new Date(), message: `${group.length} ${word} hinzugefügt.`});
            } else if (key == ActionType.DELETED) {
              return new GenericEvent({timestamp: new Date(), message: `${group.length} ${word} gelöscht.`});
            } else {
              return new GenericEvent({timestamp: new Date(), message: group.map(g => g.payload.message).join(' ,')});
            }

          });

          this.ngZone.run(() => this.snackbar.openFromComponent(EventSnackbarComponent, {
            data: messages,
            duration: 3000
          }));
        });
    });
  }
}

export interface NotificationEvent {
  payload: { timestamp: Date, message: string }

  toString();
}

export class GenericEvent implements NotificationEvent {
  constructor(public payload: { timestamp: Date, message: string }) {
  }

  public toString() {
    return `${this.payload.timestamp}: ${this.payload.message}`;
  }
}

export enum ActionType {
  UPDATED = "updated",
  ADDED = "added",
  DELETED = "deleted",
  NONE = "none"
}

export class RecordEvent implements NotificationEvent {
  constructor(public payload: { type: ActionType, message: string, timestamp: Date, record: Record }) {
  }

  public toString() {
    return `${this.payload.timestamp}, ${this.payload.record.id} ${this.payload.type}: ${this.payload.message}`;
  }
}
