import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Record, RecordService, RequiredAction} from "../../store";
import {PhysicianService} from "../physician.service";

@Component({
  selector: 'app-page-escalated',
  templateUrl: './page-escalated.component.html',
  styleUrls: ['./page-escalated.component.scss']
})
export class PageEscalatedComponent implements OnInit {

  public records: Observable<Record[]>;

  constructor(private physicianService: PhysicianService,
              private recordService: RecordService) {
    this.records = physicianService.getEscalated();
  }

  ngOnInit() {
  }

  selectRecord(record: Record) {
    this.physicianService.selectIds([record.id]);
  }

  setRequiredAction(data: { record: Record, action: RequiredAction }) {
    this.recordService.update(data.record.id, {requiredAction: data.action})
  }

  deleteRecord(record: Record) {
    this.recordService.delete(record.id)
  }

}
