import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { combineLatest, Observable } from "rxjs";
import { Category, CategoryService } from "../../shared/category-service";
import { Record, RecordService } from "../../store";
import { DocumentEditDialogComponent } from "../../shared/document-edit-dialog/document-edit-dialog.component";
import { MatDialog } from "@angular/material";
import { distinctUntilChanged, filter, map, take } from "rxjs/operators";
import { findIndex, groupBy, sortBy } from 'lodash-es';

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent implements OnInit {
  @Input() records: Observable<Record[]>;
  @Input() selectedCategory: Observable<string>;
  @Input() categories: Observable<{ [id: string]: Category }>;
  @Output() clickRecord = new EventEmitter<string>();
  @Output() selectedCategoryChange = new EventEmitter<string>();

  groupedRecords: Observable<{ category: string, records: Record[] }[]>;
  selectedIndex: Observable<number>;

  constructor(private categoryService: CategoryService,
              private recordService: RecordService,
              private dialog: MatDialog) {
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
    this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record,
      width: "635px"
    }).afterClosed().subscribe((result: Record) => {
      if (!result) {
        return;
      }
      this.recordService.update(result.id, {
        patientId: result.patientId,
        date: result.date,
        tags: result.tags,
        category: result.category
      });
    });
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
