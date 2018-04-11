import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RecordService, RequiredAction } from "../../store";

@Component({
  selector: 'app-action-bar',
  templateUrl: './action-bar.component.html',
  styleUrls: ['./action-bar.component.scss']
})
export class ActionBarComponent implements OnInit {
  @Input() recordIds: string[];
  @Output() selectAll = new EventEmitter<boolean>();

  constructor(private recordService: RecordService) { }

  ngOnInit() {
  }

  onDeleteRecord(event) {
    event.stopPropagation();
    this.recordIds.forEach(id => this.recordService.delete(id));
  }

  setRequiredAction(action: RequiredAction) {
    this.recordIds.forEach(id => this.recordService.update(id, {requiredAction: action}));
  }

  setSelection(selection: boolean) {
    this.selectAll.emit(selection);
  }
}
