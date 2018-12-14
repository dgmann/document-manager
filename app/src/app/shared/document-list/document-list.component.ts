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
import { MatSort, MatTableDataSource } from "@angular/material";
import { Router } from "@angular/router";
import { includes } from 'lodash-es';
import { DropEvent } from "ng-drag-drop";
import { Observable } from "rxjs";


import { Record, RecordService, Status } from "../../core/store";
import { CommentDialogService } from "../comment-dialog/comment-dialog.service";
import { DocumentEditDialogService } from "../document-edit-dialog/document-edit-dialog.service";


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

  displayedColumns = ['receivedAt', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();

  constructor(private recordService: RecordService,
              private router: Router,
              private commmentDialog: CommentDialogService,
              private editDialog: DocumentEditDialogService) {
  }

  ngOnInit() {
    this.dataSource.sort = this.sort;
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
    event.nativeEvent.preventDefault();
    this.appendRecord({
      source: source,
      target: event.dragData
    });
  }

  appendRecord(event) {
    this.recordService.append(event.source.id, event.target.id);
  }

  setStatus(event) {
    if (event.status == Status.ESCALATED
      || (event.status == Status.INBOX && event.record.status == Status.ESCALATED)) {
      this.commmentDialog.open(event.record).subscribe((comment: string) => {
        this.recordService.update(event.record.id, {status: event.status, comment: comment});
      });
    } else {
      this.recordService.update(event.record.id, {status: event.status});
    }
  }

  editRecord(record: Record) {
    this.editDialog.open(record).subscribe(result => {
      if (!result) {
        return;
      }
      const changes = {
        ...result.change,
        status: undefined
      };
      if (changes.patientId && changes.date && changes.category) {
        changes.status = Status.REVIEW;
      }
      this.recordService.update(result.id, changes);
    });
  }

  isSelected(selectedIds: string[], id: string) {
    if (!id) {
      return false;
    }
    return includes(selectedIds, id);
  }

  openEditor(record: Record) {
    this.router.navigate(['/editor', record.id]);
  }
}
