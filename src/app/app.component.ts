import { Component } from '@angular/core';
import { Observable } from "rxjs/Observable";


import { RecordService, Record } from "./api";
import { ReplaySubject } from "rxjs/ReplaySubject";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedRecordId = new ReplaySubject<string>();

  constructor(private recordService: RecordService) {
    this.data = recordService.all();
    this.selectedRecord = this.selectedRecordId.distinctUntilChanged().switchMap(id => this.recordService.find(id));
  }

  selectRecord(record: Record) {
    this.selectedRecordId.next(record.id);
  }

  updatePages(event) {
    this.recordService.update(event.id, {pages: event.pages})
  }
}
