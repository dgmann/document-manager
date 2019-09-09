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
import {MatSort} from '@angular/material/sort';
import {MatTableDataSource} from '@angular/material/table';
import {Router} from '@angular/router';
import {includes, without} from 'lodash-es';
import {DropEvent} from 'ng-drag-drop';


import {Record, RecordService, Status} from '../../core/store';
import {CommentDialogService} from '../comment-dialog/comment-dialog.service';
import {DocumentEditDialogService} from '../document-edit-dialog/document-edit-dialog.service';
import {NotificationService} from '@app/core';


@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  @ViewChild(MatSort, { static: false }) sort: MatSort;

  @Input() selectedIds: string[];

  @Input() set records(records: Record[]) {
    this.dataSource.data = records;
  }
  @Input() allowMultiselect = true;

  @Output() selectRecord = new EventEmitter<Record>();
  @Output() selectRecords = new EventEmitter<string[]>();

  displayedColumns = ['receivedAt', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();

  constructor(private recordService: RecordService,
              private router: Router,
              private commmentDialog: CommentDialogService,
              private notificationService: NotificationService,
              private editDialog: DocumentEditDialogService) {
  }

  ngOnInit() {
    this.dataSource.sort = this.sort;
  }

  /**
   * Set the sort after the view init since this component will
   * be able to query its view for the initialized sort.
   */
  ngAfterViewInit() {
    this.dataSource.sort = this.sort;
  }

  selectRow(record: Record, event: MouseEvent) {
    let newSelectedIds = [];
    if (this.allowMultiselect && event.getModifierState('Control')) {
      if (includes(this.selectedIds, record.id)) {
        newSelectedIds = without(this.selectedIds, record.id);
      } else {
        newSelectedIds = [...this.selectedIds, record.id];
      }
    } else {
      newSelectedIds = [record.id];
    }
    this.selectRecord.emit(record);
    this.selectRecords.emit(newSelectedIds);
  }

  deleteRecord(record: Record) {
    this.recordService.delete(record.id);
  }

  drop(source: Record, event: DropEvent) {
    event.nativeEvent.preventDefault();
    if (source.id === event.dragData.id) {
      this.notificationService.publishMessage('Anhängen an den Quellbefund nicht möglich');
      return;
    }
    this.appendRecord({
      source,
      target: event.dragData
    });
  }

  appendRecord(event) {
    this.recordService.append(event.source.id, event.target.id);
  }

  setStatus(event) {
    if (event.status === Status.ESCALATED
      || (event.status === Status.INBOX && event.record.status === Status.ESCALATED)) {
      this.commmentDialog.open(event.record).subscribe((comment: string) => {
        this.recordService.update(event.record.id, {status: event.status, comment});
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

  openEditor(record: Record) {
    this.router.navigate(['/editor', record.id]);
  }

  onDuplicateRecord(record: Record) {
    this.recordService.duplicate(record.id);
  }
}
