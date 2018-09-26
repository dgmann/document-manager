import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { flatMap, uniq } from "lodash-es";
import { Moment } from "moment";
import { BehaviorSubject, combineLatest, Observable } from "rxjs";
import { map, merge } from "rxjs/operators";
import { Patient } from "..";
import { Category } from "../../core";
import { Record } from "../../core/store/index";
import { Filter } from "../store/patient.reducer";

@Component({
  selector: 'app-record-filter',
  templateUrl: './record-filter.component.html',
  styleUrls: ['./record-filter.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RecordFilterComponent implements OnInit {
  @Input() records: Observable<Record[]>;
  @Input() patient: Observable<Patient>;
  @Output() change = new EventEmitter<Filter>();

  public tags: Observable<string[]>;
  public categories: Observable<Category[]>;
  private dateRange = new BehaviorSubject<{ from: Moment, until: Moment }>({from: null, until: null});
  private selectedTags = new BehaviorSubject<string[]>([]);
  private selectedCategories = new BehaviorSubject<Category[]>([]);

  constructor() {
  }

  ngOnInit() {
    this.tags = this.patient.pipe(
      map(patient => patient && patient.tags || []),
      merge(this.records.pipe(map(records => uniq(flatMap(records, r => r.tags)))))
    );
    this.categories = this.patient.pipe(
      map(patient => patient && patient.categories || [])
    );


    combineLatest(
      this.selectedTags,
      this.selectedCategories,
      this.dateRange
    ).subscribe(([tags, categories, dateRange]) => this.change.emit({
      categoryIds: categories.map(c => c.id),
      tags: tags,
      from: dateRange.from,
      until: dateRange.until
    }));
  }

  onSelectTags(tags: string[]) {
    this.selectedTags.next(tags);
  }

  onSelectCategories(categories: Category[]) {
    this.selectedCategories.next(categories);
  }

  onDateFilterChange(value: { from: Moment, until: Moment }) {
    this.dateRange.next(value);
  }

}
