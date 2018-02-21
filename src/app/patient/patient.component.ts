import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {of} from "rxjs/observable/of";
import {Record, RecordService} from "../store";

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss']
})
export class PatientComponent implements OnInit {
  public patientId: Observable<string>;
  public records: Observable<Record[]>;

  constructor(private recordService: RecordService) {
    this.patientId = of("3");
    this.records = recordService.all();
  }

  ngOnInit() {
  }

}
