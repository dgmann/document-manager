import {Injectable} from "@angular/core";
import {WebSocketSubject} from "rxjs/observable/dom/WebSocketSubject";

@Injectable()
export class WebsocketService {
  public create(url: string) {
    return new WebSocketSubject<NotificationMessage>(url);
  }
}

export interface NotificationMessage {
  type: NotificationMessageType;
  timestamp: Date;
  data: any;
}

export enum NotificationMessageType {
  Created = "CREATE",
  Updated = "UPDATE",
  Deleted = "DELETE"
}
