import {Component, Input, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Record} from "../../store";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss']
})
export class MultiRecordListComponent implements OnInit {
  @Input('records') records: Observable<Record[]>;

  constructor() {
  }

  ngOnInit() {
  }

}
