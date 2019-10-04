import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Record, Status} from '../../core/store/index';

@Component({
  selector: 'app-action-menu',
  templateUrl: './action-menu.component.html',
  styleUrls: ['./action-menu.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ActionMenuComponent implements OnInit {
  @Input() record: Record;
  @Output() deleteRecord = new EventEmitter<Record>();
  @Output() changeStatus = new EventEmitter<{ record: Record, status: Status }>();
  @Output() editRecord = new EventEmitter<Record>();
  @Output() export = new EventEmitter<Record>();
  @Output() duplicateRecord = new EventEmitter<Record>();

  status = Status;

  constructor() {
  }

  ngOnInit() {
  }

  onDeleteRecord(event, row: Record) {
    this.deleteRecord.emit(row);
    event.stopPropagation();
  }

  setStatus(record: Record, action: Status) {
    this.changeStatus.emit({record, status: action});
  }

  onEditRecord(record: Record) {
    this.editRecord.emit(record);
  }

  onExportAsPdf(record: Record) {
    this.export.emit(record);
  }

  onDuplicateRecord(record: Record) {
    this.duplicateRecord.emit(record);
  }
}
