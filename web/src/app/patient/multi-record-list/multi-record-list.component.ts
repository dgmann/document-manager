import {ChangeDetectionStrategy, Component, EventEmitter, Input, Output} from '@angular/core';
import {findIndex, groupBy, sortBy} from 'lodash-es';
import {Category} from '@app/core/categories';
import {Record} from '@app/core/records';
import {DocumentEditDialogService, EditResult, MessageBoxService} from '../../shared';
import {Patient} from '@app/patient';

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent {
  patientName: string;

  @Output() clickRecord = new EventEmitter<string>();
  @Output() updateRecord = new EventEmitter<EditResult>();
  @Output() deleteRecord = new EventEmitter<Record>();
  @Output() openInEditor = new EventEmitter<Record>();
  @Output() duplicate = new EventEmitter<Record>();

  @Input() set patient(patient: Patient) {
    if (patient) {
      this.patientName = `${patient.lastName}, ${patient.firstName}`;
    }
  }
  @Input() records: Record[];

  @Input() categories: { [id: string]: Category };
  @Input() selectedRecordId: string;

  constructor(private dialog: DocumentEditDialogService, private messageBox: MessageBoxService) {
  }

  onRecordClicked(id: string) {
    this.clickRecord.emit(id);
  }

  onEdit(record: Record) {
    this.dialog.open(record).subscribe(result => this.updateRecord.emit(result));
  }

  onDelete(record: Record) {
    this.messageBox.open('Löschen', 'Wollen sie diesen Befund löschen?').subscribe(yes => {
      if (yes) {
        this.deleteRecord.emit(record);
      }
    });
  }

  onOpenInEditor(record: Record) {
    this.openInEditor.emit(record);
  }

  onDuplicateRecord(record: Record) {
    this.duplicate.emit(record);
  }
}
