import {ChangeDetectionStrategy, Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Moment} from "moment";
import {BehaviorSubject, combineLatest} from "rxjs";
import {Filter} from "../store/patient.reducer";

@Component({
  selector: 'app-date-range-selector',
  templateUrl: './date-range-selector.component.html',
  styleUrls: ['./date-range-selector.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DateRangeSelectorComponent implements OnInit {
  @Output() change = new EventEmitter<Filter>();
  untilValue: Date = null;
  fromValue: Date = null;
  private from = new BehaviorSubject<Moment>(null);
  private until = new BehaviorSubject<Moment>(null);

  constructor() {
  }

  ngOnInit() {
    combineLatest(
      this.from,
      this.until
    ).subscribe(([from, until]) => this.change.emit({
      from: from,
      until: until
    }));
  }

  setFrom(value: Moment) {
    this.from.next(value);
  }

  setUntil(value: Moment) {
    this.until.next(value);
  }

  clearFrom() {
    this.fromValue = null;
    this.from.next(null);
  }

  clearUntil() {
    this.untilValue = null;
    this.until.next(null);
  }
}
