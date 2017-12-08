import { Component, OnInit, ViewChild, AfterViewInit, Input } from '@angular/core';
import { MatSort, MatTableDataSource } from "@angular/material";
import { Record } from "../core/record";

@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss']
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  displayedColumns = ['id', 'date', 'type', 'comment'];
  dataSource: MatTableDataSource<Record>;

  @ViewChild(MatSort) sort: MatSort;
  @Input('data') data: Record[];

  constructor() { }

  ngOnInit() {
    this.dataSource = new MatTableDataSource(this.data);
  }

  /**
   * Set the sort after the view init since this component will
   * be able to query its view for the initialized sort.
   */
  ngAfterViewInit() {
    this.dataSource.sort = this.sort;
  }
}

export interface Element {
  name: string;
  position: number;
  weight: number;
  symbol: string;
}
