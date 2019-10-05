import {Component, OnInit} from '@angular/core';
import {HistoryService} from '@app/history/history-service';
import {Record} from '@app/core/store';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';

@Component({
  selector: 'app-history',
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.scss']
})
export class HistoryComponent implements OnInit {
  records$: Observable<Record[]>;
  selectedIds$: Observable<string[]>;
  selectedRecord$: Observable<Record>;

  constructor(private historyService: HistoryService) {
    this.historyService.next();
    this.records$ = this.historyService.records$;
    this.records$.subscribe(r => console.log(r));
    this.selectedIds$ = this.historyService.selectedId$.pipe(map(id => id && [id] || []));
    this.selectedRecord$ = this.historyService.selectedRecord$;
  }

  ngOnInit() {

  }

  onSelectRecords(ids: string[]) {
    this.historyService.selectRecord(ids[0]);
  }

}
