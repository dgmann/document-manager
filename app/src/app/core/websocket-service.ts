import {Injectable} from '@angular/core';
import {webSocket} from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  public create(url: string) {
    return webSocket<NotificationMessage>(url);
  }
}

export interface NotificationMessage {
  type: NotificationMessageType;
  timestamp: Date;
  id: string;
  topic: NotificationTopic;
}

export enum NotificationMessageType {
  Created = 'CREATE',
  Updated = 'UPDATE',
  Deleted = 'DELETE'
}

export enum NotificationTopic {
  Records = 'records'
}
