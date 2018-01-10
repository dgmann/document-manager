import { AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatSort, MatTableDataSource } from "@angular/material";
import { Observable } from "rxjs/Observable";


import { Record } from "../api";

@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss']
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  displayedColumns = ['date', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();
  selectedRecordId = "";

  @ViewChild(MatSort) sort: MatSort;
  @Input('data') data: Observable<Record[]>;
  @Output('recordClicked') recordClicked = new EventEmitter<Record>();

  constructor() { }

  ngOnInit() {
    this.data.subscribe(data => this.dataSource.data = data);
  }

  /**
   * Set the sort after the view init since this component will
   * be able to query its view for the initialized sort.
   */
  ngAfterViewInit() {
    this.dataSource.sort = this.sort;
  }

  selectRow(row: Record) {
    this.recordClicked.emit(row);
    this.selectedRecordId = row.id
  }
}
