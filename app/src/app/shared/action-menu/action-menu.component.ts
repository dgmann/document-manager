import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Record, RequiredAction} from "../../store";

@Component({
  selector: 'app-action-menu',
  templateUrl: './action-menu.component.html',
  styleUrls: ['./action-menu.component.scss']
})
export class ActionMenuComponent implements OnInit {
  @Input() record: Record;
  @Output() deleteRecord = new EventEmitter<Record>();
  @Output() changeRequiredAction = new EventEmitter<{ record: Record, action: RequiredAction }>();

  constructor() {
  }

  ngOnInit() {
  }

  onDeleteRecord(event, row: Record) {
    this.deleteRecord.emit(row);
    event.stopPropagation();
  }

  setRequiredAction(record: Record, action: RequiredAction) {
    this.changeRequiredAction.emit({record: record, action: action});
  }

}
