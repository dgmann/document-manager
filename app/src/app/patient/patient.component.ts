import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {of} from "rxjs/observable/of";
import {map, switchMap} from "rxjs/operators";
import {Patient, PatientService} from "../shared";
import {Category} from "../shared/category-service";
import {Record, RecordService} from "../store";

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss']
})
export class PatientComponent implements OnInit {
  public patientId: Observable<string>;
  public records: Observable<Record[]>;
  public patient: Observable<Patient>;
  public tags: Observable<string[]>;
  public categories: Observable<Category[]>;

  constructor(private recordService: RecordService, private patientService: PatientService) {
    this.patientId = of("3");
    this.records = recordService.all();
    this.patient = this.patientId.pipe(
      switchMap(id => patientService.findById(id))
    );
    this.tags = this.patient.pipe(
      map(patient => patient.tags)
    );
    this.categories = this.patient.pipe(
      map(patient => patient.categories)
    );
  }

  ngOnInit() {
  }

}
