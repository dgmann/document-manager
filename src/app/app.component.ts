import { Component } from '@angular/core';
import { Observable } from "rxjs/Observable";


import { RecordService, Record } from "./api";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  data: Observable<Record[]>;
  selectedRecord: Record;

  constructor(private recordService: RecordService) {
    this.data = recordService.all();
  }

  selectRecord(record: Record) {
    this.selectedRecord = record;
  }
}
