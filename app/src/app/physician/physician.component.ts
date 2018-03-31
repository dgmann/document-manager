import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {filter, map} from "rxjs/operators";
import {Record, RecordService, RequiredAction} from "../store";
import {PhysicianService} from "./physician.service";

@Component({
  selector: 'app-physician',
  templateUrl: './physician.component.html',
  styleUrls: ['./physician.component.scss']
})
export class PhysicianComponent implements OnInit {

  private selectedRecord: Observable<Record>;

  constructor(private physicianService: PhysicianService,
              private recordService: RecordService) {
    this.selectedRecord = physicianService.getSelectedRecords().pipe(
      filter(records => records.length > 0),
      map(records => records[0])
    );
    this.physicianService.load(RequiredAction.REVIEW);
    this.physicianService.load(RequiredAction.ESCALATED);
    this.physicianService.load(RequiredAction.OTHER);
  }

  ngOnInit() {
  }

}
