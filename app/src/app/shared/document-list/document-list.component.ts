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


import {Record, RecordService, Status} from '@app/core/records';
import {CommentDialogService} from '../comment-dialog/comment-dialog.service';
import {DocumentEditDialogService} from '../document-edit-dialog/document-edit-dialog.service';
import {NotificationService} from '@app/core/notifications';


@Component({
  selector: 'app-document-list',
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentListComponent implements OnInit, AfterViewInit {
  static readonly DRAG_TYPE = 'recordid';

  @ViewChild(MatSort) sort: MatSort;

  @Input() selectedIds: string[];

  @Input() set records(records: Record[]) {
    this.dataSource.data = records;
  }

  get records() {
    return this.dataSource.sortData(this.dataSource.filteredData, this.dataSource.sort);
  }

  @Input() allowMultiselect = true;

  @Output() selectRecord = new EventEmitter<Record>();
  @Output() selectRecords = new EventEmitter<string[]>();

  displayedColumns = ['receivedAt', 'sender', 'numpages', 'comment', 'actions'];
  dataSource = new MatTableDataSource<Record>();
  isDragging = false;

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
    const selectedRecordId = record.id;
    if (this.allowMultiselect && event.getModifierState('Control')) {
      if (includes(this.selectedIds, record.id)) {
        newSelectedIds = without(this.selectedIds, selectedRecordId);
      } else {
        newSelectedIds = [...this.selectedIds, selectedRecordId];
      }
    } else if (this.allowMultiselect && event.getModifierState('Shift')) {
      if (this.selectedIds.length > 0) {
        const startIndex = this.records.findIndex(r => r.id === this.selectedIds[0]);
        const endIndex = this.records.findIndex(r => r.id === selectedRecordId);
        const boundaries = [startIndex, endIndex].sort();
        newSelectedIds = this.records.slice(boundaries[0], boundaries[1] + 1).map(r => r.id);
      } else {
        newSelectedIds = [...this.selectedIds, selectedRecordId];
      }
    } else {
      newSelectedIds = [selectedRecordId];
    }
    const target = event.currentTarget as HTMLElement;
    target.focus();
    this.selectRecord.emit(record);
    this.selectRecords.emit(newSelectedIds);
  }

  deleteRecord(record: Record) {
    this.recordService.delete(record.id);
  }

  onDragStart(event: DragEvent, record: Record) {
    event.dataTransfer.setData(DocumentListComponent.DRAG_TYPE, record.id);
    this.isDragging = true;
  }

  setDragOver(event: DragEvent, dragover: boolean) {
    const target = event.currentTarget as HTMLElement;
    target.setAttribute('dragover', dragover + '');
  }

  onDrop(event: DragEvent) {
    const target = event.currentTarget as HTMLElement;
    if (!target || (target && target.tagName !== 'MAT-ROW')) {
      return;
    }

    event.preventDefault();
    const sourceRecordId = event.dataTransfer.getData('recordId');
    const targetRecordId = target.getAttribute('recordid');

    if (sourceRecordId === targetRecordId) {
      this.notificationService.publishMessage('Anhängen an den Quellbefund nicht möglich');
      return;
    }
    this.recordService.append(sourceRecordId, targetRecordId);
  }

  onDragEnd(event: DragEvent) {
    this.isDragging = false;
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

  containsDragType(event: DragEvent): boolean {
    if (event.dataTransfer.types) {
      for (const type of event.dataTransfer.types) {
        if (type === DocumentListComponent.DRAG_TYPE) {
          return true;
        }
      }
    }
    return false;
  }
}
