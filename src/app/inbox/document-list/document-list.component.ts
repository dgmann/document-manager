import { AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatSort, MatTableDataSource } from "@angular/material";
import { Observable } from "rxjs/Observable";


import { Record } from "../../store/index";
import { DropEvent } from "ng-drag-drop";

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
  @Output('recordDelete') recordDelete = new EventEmitter<Record>();
  @Output('recordClick') recordClick = new EventEmitter<Record>();
  @Output('recordDbClick') recordDbClick = new EventEmitter<Record>();
  @Output('recordDrop') recordDrop = new EventEmitter<{ source: Record, target: Record }>();

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
    this.recordClick.emit(row);
    this.selectedRecordId = row.id
  }

  rowDoubleClick(record: Record) {
    this.recordDbClick.emit(record);
  }

  deleteRecord(event, row: Record) {
    this.recordDelete.emit(row);
    event.stopPropagation();
  }

  drop(source: Record, event: DropEvent) {
    this.recordDrop.emit({
      source: source,
      target: event.dragData
    });
  }
}