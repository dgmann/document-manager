import {ChangeDetectionStrategy, Component, EventEmitter, Output} from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import {Moment} from 'moment';
import { distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-date-range-selector',
  templateUrl: './date-range-selector.component.html',
  styleUrls: ['./date-range-selector.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DateRangeSelectorComponent {
  @Output() dateRangeChange = new EventEmitter<{ from: Moment, until: Moment }>();
  range = new FormGroup({
    start: new FormControl(),
    end: new FormControl()
  });

  constructor() {
  }

  onDateChange() {
    this.dateRangeChange.emit({from: this.range.value.start, until: this.range.value.end});
  }
}
