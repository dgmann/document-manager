import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild
} from '@angular/core';
import {MatDialog, MatSort, MatTableDataSource} from "@angular/material";
import {Router} from "@angular/router";
import {includes} from 'lodash-es';
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs";


import {Record, RecordService, Status} from "../../store";
import {DocumentEditDialogComponent} from "../document-edit-dialog/document-edit-dialog.component";

@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  @ViewChild(MatSort) sort: MatSort;

  @Input() selectedIds: Observable<string[]>;
  @Input() records: Observable<Record[]>;
  @Output() selectRecord = new EventEmitter<Record>();

  displayedColumns = ['date', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();

  constructor(private recordService: RecordService,
              private router: Router,
              private dialog: MatDialog) {
  }

  ngOnInit() {
    this.records.subscribe(data => this.dataSource.data = data);
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

  setStatus(event) {
    this.recordService.update(event.record.id, {status: event.status});
  }

  editRecord(record: Record) {
    this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record,
      width: "635px"
    }).afterClosed().subscribe((result: Record) => {
      if (!result) {
        return;
      }
      const changes = {
        patientId: result.patientId,
        date: result.date,
        tags: result.tags,
        categoryId: result.categoryId,
        status: undefined
      };
      if (changes.patientId && changes.date && changes.categoryId) {
        changes.status = Status.REVIEW;
      }
      this.recordService.update(result.id, changes);
    });
  }

  isSelected(selectedIds: string[], id: string) {
    return includes(selectedIds, id);
  }

  openEditor(record: Record) {
    this.router.navigate(['/editor', record.id]);
  }
}
