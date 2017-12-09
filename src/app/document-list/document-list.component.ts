import { Component, OnInit, ViewChild, AfterViewInit, Input } from '@angular/core';
import { MatSort, MatTableDataSource } from "@angular/material";
import { Observable } from "rxjs/Observable";


import { Record } from "../api";

@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss']
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  displayedColumns = ['id', 'date', 'sender', 'comment'];
  dataSource: MatTableDataSource<Record>;

  @ViewChild(MatSort) sort: MatSort;
  @Input('data') data: Observable<Record[]>;

  constructor() { }

  ngOnInit() {
    this.data.subscribe(data => this.dataSource = new MatTableDataSource(data));
  }

  /**
   * Set the sort after the view init since this component will
   * be able to query its view for the initialized sort.
   */
  ngAfterViewInit() {
    this.dataSource.sort = this.sort;
  }
}
