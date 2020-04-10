import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {select, Store} from '@ngrx/store';
import {map} from 'rxjs/operators';
import {State} from '../store';
import {DeleteRecord, LoadRecords, UpdatePages, UpdateRecord} from './record.actions';
import {PageUpdate, Record} from './record.model';
import {selectAllRecords, selectInvalidIds, selectIsLoading, selectRecordEntities} from './record.selectors';
import {Observable} from 'rxjs';
import {ActionType, GenericEvent, NotificationService, RecordEvent} from '../notifications/notification-service';
import {ConfigService} from '@app/core/config';

@Injectable({
  providedIn: 'root'
})
export class RecordService {
  public records: Observable<Record[]>;
  public invalidIds: Observable<string[]>;
  public isLoading$: Observable<boolean>;

  constructor(private store: Store<State>,
              private http: HttpClient,
              private notifications: NotificationService,
              private config: ConfigService) {
    this.records = this.store.pipe(select(selectAllRecords));
    this.invalidIds = this.store.pipe(select(selectInvalidIds));
    this.isLoading$ = this.store.pipe(select(selectIsLoading));
  }

  public load(query: { [param: string]: string }) {
    this.store.dispatch(new LoadRecords({query}));
  }

  public find(id: string) {
    return this.store.pipe(select(selectRecordEntities), map(entities => entities[id]));
  }

  public delete(id: string) {
    this.store.dispatch(new DeleteRecord({id}));
  }

  public update(id: string, changes: any) {
    this.store.dispatch(new UpdateRecord({record: {id, changes}}));
  }

  public updatePages(id: string, pages: PageUpdate[]) {
    this.store.dispatch(new UpdatePages({id, updates: pages}));
  }

  public upload(pdf) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: 'PDF wird hochgeladen...'
    }));

    const formData = new FormData();
    formData.append('pdf', pdf);
    formData.append('sender', 'Client');
    this.http.post<Record>(this.config.getApiUrl() + '/records', formData)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: 'PDF hochgeladen',
        id: record.id
    })));
  }

  public append(sourceId: string, targetId: string) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: 'PDF wird angehängt...'
    }));
    this.http.post<Record>(`${this.config.getApiUrl()}/records/${targetId}/append/${sourceId}`, null)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: 'PDF angehängt',
        id: record.id
      })));
  }

  public reset(id: string) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: 'Befund wird zurückgesetzt...'
    }));
    this.http.put<Record>(`${this.config.getApiUrl()}/records/${id}/reset`, null)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: 'Befund zurückgesetzt',
        id: record.id
    })));
  }

  public duplicate(id: string) {
    this.notifications.publish(new GenericEvent({
      timestamp: new Date(),
      message: 'Befund wird dupliziert...'
    }));
    this.http.post<Record>(`${this.config.getApiUrl()}/records/${id}/duplicate`, null)
      .subscribe(record => this.notifications.publish(new RecordEvent({
        type: ActionType.NONE,
        timestamp: new Date(),
        message: 'Befund dupliziert',
        id: record.id
      })));
  }

  public createPDFLink(ids: string[]) {
    const url = new URL(`${this.config.getApiUrl()}/export`);
    ids.forEach(id => url.searchParams.append('id', id));
    return url.href;
  }
}
