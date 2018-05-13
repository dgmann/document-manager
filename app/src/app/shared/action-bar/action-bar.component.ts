import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {RecordService, Status} from "../../store";

@Component({
  selector: 'app-action-bar',
  templateUrl: './action-bar.component.html',
  styleUrls: ['./action-bar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ActionBarComponent implements OnInit {
  @Input() recordIds: string[];
  @Output() selectAll = new EventEmitter<boolean>();

  status = Status;

  constructor(private recordService: RecordService) { }

  ngOnInit() {
  }

  onDeleteRecord(event) {
    event.stopPropagation();
    this.recordIds.forEach(id => this.recordService.delete(id));
  }

  setStatus(action: Status) {
    this.recordIds.forEach(id => this.recordService.update(id, {status: action}));
  }

  setSelection(selection: boolean) {
    this.selectAll.emit(selection);
  }
}
