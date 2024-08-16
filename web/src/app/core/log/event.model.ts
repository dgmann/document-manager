import {Record} from '../records';

export interface Event {
  payload: { timestamp: Date, message: string };

  toString();
}

export class GenericEvent implements Event {
  constructor(public payload: { timestamp: Date, message: string }) {
  }

  public toString() {
    return `${this.payload.timestamp}: ${this.payload.message}`;
  }
}

export enum ActionType {
  UPDATED = 'updated',
  ADDED = 'added',
  DELETED = 'deleted',
  NONE = 'none'
}

export class RecordEvent implements Event {
  constructor(public payload: { type: ActionType, message: string, timestamp: Date, record: Record }) {
  }

  public toString() {
    return `${this.payload.timestamp}, ${this.payload.record.id} ${this.payload.type}: ${this.payload.message}`;
  }
}
