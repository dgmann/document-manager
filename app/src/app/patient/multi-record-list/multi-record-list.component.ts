import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {findIndex, groupBy, sortBy} from 'lodash-es'
import {combineLatest, Observable} from "rxjs";
import {distinctUntilChanged, filter, map, take, withLatestFrom} from "rxjs/operators";
import {Category} from "../../core";
import {Record} from "../../core/store";
import {DocumentEditDialogService, EditResult, MessageBoxService} from "../../shared";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent implements OnInit {
  @Input() records: Observable<Record[]>;
  @Input() selectedCategory: Observable<string>;
  categoryMap: { [id: string]: Category } = {};
  @Output() clickRecord = new EventEmitter<string>();
  @Output() selectedCategoryChange = new EventEmitter<string>();
  @Output() updateRecord = new EventEmitter<EditResult>();
  @Output() deleteRecord = new EventEmitter<Record>();
  @Output() openInEditor = new EventEmitter<Record>();
  allRecords: Observable<Record[]>;
  groupedRecords: Observable<{ category: string, records: Record[] }[]>;

  @Input() set categories(value: { [id: string]: Category }) {
    this.categoryMap = {
      ...value,
      all: {id: "all", name: "Alle"}
    };
  }
  selectedIndex: Observable<number>;

  constructor(private dialog: DocumentEditDialogService, private messageBox: MessageBoxService) {
  }

  ngOnInit() {
    this.allRecords = this.records.pipe(
      map(records => sortBy(records, ['category', 'date']))
    );

    this.groupedRecords = this.records.pipe(
      map(records => groupBy(records, 'category')),
      map(grouped => Object.entries(grouped)
        .map(entry => ({category: entry[0], records: sortBy(entry[1], ['date'])}))),
      map(grouped => sortBy(grouped, ['category'])),
      withLatestFrom(this.records),
      map(([grouped, records]) => [...grouped, {category: 'all', records: sortBy(records, ['date'])}])
    );
    this.selectedIndex = combineLatest(this.groupedRecords, this.selectedCategory).pipe(
      map(([categories, selectedCategory]) => findIndex(categories, {category: selectedCategory})),
      distinctUntilChanged()
    );
  }

  onRecordClicked(id: string) {
    this.clickRecord.emit(id);
  }

  edit(record: Record) {
    this.dialog.open(record).subscribe(result => this.updateRecord.emit(result));
  }

  delete(record: Record) {
    this.messageBox.open("Löschen", "Wollen sie diesen Befund löschen?").subscribe(yes => {
      if (yes) {
        this.deleteRecord.emit(record);
      }
    });
  }

  onOpenInEditor(record: Record) {
    this.openInEditor.emit(record);
  }

  onSelectedIndexChange(index: number) {
    this.groupedRecords
      .pipe(
        take(1),
        map(groups => groups[index] && groups[index].category || null),
        filter(c => !!c)
      ).subscribe(category => this.selectedCategoryChange.emit(category));
  }
}
