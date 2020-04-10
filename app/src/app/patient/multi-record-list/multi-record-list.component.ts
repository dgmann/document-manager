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
  categoryMap: { [id: string]: Category } = {};
  groupedRecords: { category: string, records: Record[] }[];
  selectedIndex: number;

  @Output() clickRecord = new EventEmitter<string>();
  @Output() selectedCategoryChange = new EventEmitter<string>();
  @Output() updateRecord = new EventEmitter<EditResult>();
  @Output() deleteRecord = new EventEmitter<Record>();
  @Output() openInEditor = new EventEmitter<Record>();
  @Output() duplicate = new EventEmitter<Record>();
  selectedCategory: string;

  @Input() set patient(patient: Patient) {
    if (patient) {
      this.patientName = `${patient.lastName}, ${patient.firstName}`;
    }
  }
  @Input() set records(records: Record[]) {
    const grouped = groupBy(records, 'category');
    const groupedRecords = Object.entries(grouped)
      .map(entry => ({category: entry[0], records: sortBy(entry[1], ['date'])}));
    const sorted = sortBy(groupedRecords, ['category']);
    this.groupedRecords = [...sorted, {category: 'all', records: sortBy(records, ['date'])}];
    this.setSelectedIndex(this.groupedRecords, this.selectedCategory);
  }

  @Input('selectedCategory') set setSelectedCategory(selectedCategory: string) {
    this.selectedCategory = selectedCategory;
    this.setSelectedIndex(this.groupedRecords, this.selectedCategory);
  }

  @Input() set categories(value: { [id: string]: Category }) {
    this.categoryMap = {
      ...value,
      all: {id: 'all', name: 'Alle'}
    };
  }

  constructor(private dialog: DocumentEditDialogService, private messageBox: MessageBoxService) {
  }

  setSelectedIndex(groupedRecords: { category: string, records: Record[] }[], selectedCategory: string) {
    this.selectedIndex = findIndex(groupedRecords, {category: selectedCategory});
  }

  onRecordClicked(id: string) {
    this.clickRecord.emit(id);
  }

  edit(record: Record) {
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

  onSelectedIndexChange(index: number) {
    const category = this.groupedRecords[index] && this.groupedRecords[index].category || null;
    if (category) {
      this.selectedCategoryChange.emit(category);
    }
  }
}
