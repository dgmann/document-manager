import {AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {MatDialog, MatSort, MatTableDataSource} from "@angular/material";
import {includes} from 'lodash-es';
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs/Observable";


import {Record, RecordService} from "../../store";
import {DocumentEditDialogComponent} from "../document-edit-dialog/document-edit-dialog.component";

@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss']
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  @ViewChild(MatSort) sort: MatSort;

  @Input() selectedIds: Observable<string[]>;
  @Input() records: Observable<Record[]>;
  @Output() selectRecord = new EventEmitter<Record>();

  displayedColumns = ['date', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();
  selectedRecordIds = [];

  constructor(private recordService: RecordService,
              private dialog: MatDialog) {
  }

  ngOnInit() {
    this.records.subscribe(data => this.dataSource.data = data);
    this.selectedIds.subscribe(ids => this.selectedRecordIds = ids);
  }

  /**
   * Set the sort after the view init since this component will
   * be able to query its view for the initialized sort.
   */
  ngAfterViewInit() {
    this.dataSource.sort = this.sort;
  }

  selectRow(row: Record) {
    this.selectRecord.emit(row);
  }

  deleteRecord(record: Record) {
    this.recordService.delete(record.id);
  }

  drop(source: Record, event: DropEvent) {
    this.appendRecord({
      source: source,
      target: event.dragData
    });
  }

  appendRecord(event) {
    this.recordService.append(event.source.id, event.target.id);
  }

  setRequiredAction(event) {
    this.recordService.update(event.record.id, {requiredAction: event.action});
  }

  editRecord(record: Record) {
    this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record
    }).afterClosed().subscribe((result: Record) => {
      if (!result) {
        return;
      }
      this.recordService.update(result.id, {
        patientId: result.patientId,
        date: result.date,
        tags: result.tags,
        categoryId: result.categoryId
      });
    });
  }

  isSelected(id: string) {
    return includes(this.selectedRecordIds, id);
  }
}
