import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-three-column-panel',
  templateUrl: './three-column-panel.component.html',
  styleUrls: ['./three-column-panel.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ThreeColumnPanelComponent implements OnInit {
  showLeftColumnValue = true;
  @Output() showLeftColumnChange = new EventEmitter<boolean>();
  showRightColumnValue = false;
  @Output() showRightColumnChange = new EventEmitter<boolean>();

  constructor() {
  }

  @Input()
  get showLeftColumn() {
    return this.showLeftColumnValue;
  }

  set showLeftColumn(val) {
    this.showLeftColumnValue = val;
    this.showLeftColumnChange.emit(val);
  }

  @Input()
  get showRightColumn() {
    return this.showRightColumnValue;
  }

  set showRightColumn(val) {
    this.showRightColumnValue = val;
    this.showRightColumnChange.emit(val);
  }

  ngOnInit() {
  }

  openLeftColumn() {
    this.showLeftColumn = true;
    this.showRightColumn = false;
  }
}
