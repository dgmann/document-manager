import { Injectable } from "@angular/core";
import { CategoryService } from "../../core";
import { Record, RecordService } from "../../core/store";
import { combineLatest, Observable } from "rxjs";
import { distinctUntilChanged, map } from "rxjs/operators";
import { findIndex, groupBy, sortBy } from "lodash-es";

@Injectable({
  providedIn: "root"
})
export class MultiRecordListService {
  groupedRecords: Observable<{ category: string, records: Record[] }[]>;
  selectedIndex: Observable<number>;

  constructor(private categoryService: CategoryService,
              private recordService: RecordService) {
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
}
