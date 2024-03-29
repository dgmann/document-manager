import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Record} from '@app/core/records';
import {map, withLatestFrom} from 'rxjs/operators';
import {BehaviorSubject, combineLatest, Observable} from 'rxjs';
import {ConfigService} from '@app/core/config';

@Injectable()
export class HistoryService {
  selectedId$ = new BehaviorSubject<string>(null);
  selectedRecord$: Observable<Record>;
  records$ = new BehaviorSubject<Record[]>([]);
  skip = 0;
  limit = 20;

  constructor(private http: HttpClient, private config: ConfigService) {
    this.selectedRecord$ = combineLatest([this.selectedId$, this.records$])
      .pipe(map(([id, records]) => records.find(record => record.id === id)));
  }

  public next() {
    this.get(this.skip, this.limit).pipe(withLatestFrom(this.records$)).subscribe(([records, oldRecords]) => {
      const newRecords = [...records, ...oldRecords];
      this.records$.next(newRecords);
    });
    this.skip += this.limit;
  }

  public selectRecord(id: string) {
    this.selectedId$.next(id);
  }

  private get(skip: number, limit: number) {
    const params = {
      sort: '-updatedAt',
      skip: skip.toString(),
      limit: limit.toString()
    };
    return this.http.get<Record[]>(`${this.config.getApiUrl()}/records`, {params});
  }
}
