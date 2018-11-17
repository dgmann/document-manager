import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { combineLatest, Observable } from "rxjs";
import { Category } from "../../core";
import { Record } from "../../core/store";
import { distinctUntilChanged, filter, map, take } from "rxjs/operators";
import { findIndex, groupBy, sortBy } from 'lodash-es'
import { DocumentEditDialogService, EditResult } from "../../shared";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent implements OnInit {
  @Input() records: Observable<Record[]>;
  @Input() selectedCategory: Observable<string>;
  @Input() categories: { [id: string]: Category };
  @Output() clickRecord = new EventEmitter<string>();
  @Output() selectedCategoryChange = new EventEmitter<string>();
  @Output() updateRecord = new EventEmitter<EditResult>();

  groupedRecords: Observable<{ category: string, records: Record[] }[]>;
  selectedIndex: Observable<number>;

  constructor(private dialog: DocumentEditDialogService) {
  }

  ngOnInit() {
    this.groupedRecords = this.records.pipe(
      map(records => groupBy(records, 'category')),
      map(grouped => Object.entries(grouped)
        .map(entry => ({category: entry[0], records: entry[1]}))),
      map(grouped => sortBy(grouped, ['category']))
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

  onSelectedIndexChange(index: number) {
    this.groupedRecords
      .pipe(
        take(1),
        map(groups => groups[index].category),
        filter(c => !!c)
      ).subscribe(category => this.selectedCategoryChange.emit(category));
  }
}
