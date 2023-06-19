import {ChangeDetectionStrategy, Component, EventEmitter, Output} from '@angular/core';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';

@Component({
  selector: 'app-date-range-selector',
  templateUrl: './date-range-selector.component.html',
  styleUrls: ['./date-range-selector.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DateRangeSelectorComponent {
  @Output() dateRangeChange = new EventEmitter<{ from: Date, until: Date }>();
  range = new UntypedFormGroup({
    start: new UntypedFormControl(),
    end: new UntypedFormControl()
  });

  constructor() {
  }

  onDateChange() {
    this.dateRangeChange.emit({from: this.range.value.start, until: this.range.value.end});
  }
}
