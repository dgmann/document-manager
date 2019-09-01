import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Moment} from 'moment';
import {BehaviorSubject, combineLatest, Observable} from 'rxjs';
import {Patient} from '..';
import {Category} from '@app/core';
import {Filter} from '../store/patient.reducer';

@Component({
  selector: 'app-record-filter',
  templateUrl: './record-filter.component.html',
  styleUrls: ['./record-filter.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RecordFilterComponent implements OnInit {
  @Input() patient: Observable<Patient>;
  @Input() categories: Observable<Category[]>;
  @Input() tags: Observable<string[]>;
  @Output() change = new EventEmitter<Filter>();

  private dateRange = new BehaviorSubject<{ from: Moment, until: Moment }>({from: null, until: null});
  private selectedTags = new BehaviorSubject<string[]>([]);
  private selectedCategories = new BehaviorSubject<Category[]>([]);

  constructor() {
  }

  ngOnInit() {
    combineLatest(
      this.selectedTags,
      this.selectedCategories,
      this.dateRange
    ).subscribe(([tags, categories, dateRange]) => this.change.emit({
      categoryIds: categories.map(c => c.id),
      tags,
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
