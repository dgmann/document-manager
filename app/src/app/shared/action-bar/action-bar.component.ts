import { ChangeDetectionStrategy, Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Status } from "../../core/store/index";

@Component({
  selector: 'app-action-bar',
  templateUrl: './action-bar.component.html',
  styleUrls: ['./action-bar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ActionBarComponent implements OnInit {
  @Output() selectAll = new EventEmitter<boolean>();
  @Output() delete = new EventEmitter<void>();
  @Output() changeStatus = new EventEmitter<Status>();

  status = Status;

  ngOnInit() {
  }

  onDeleteRecord(event) {
    event.stopPropagation();
    this.delete.emit();
  }

  setStatus(action: Status) {
    this.changeStatus.emit(action);
  }

  setSelection(selection: boolean) {
    this.selectAll.emit(selection);
  }
}
