import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {map} from "rxjs/operators";
import {Patient} from "../shared";
import {Category} from "../shared/category-service";
import {Record, RecordService} from "../store";
import {PatientService} from "./patient.service";

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
    patientService.selectPatient('3');
    this.records = patientService.getPatientRecords();
    this.patient = patientService.getSelectedPatient();
    this.patientId = this.patient.pipe(
      map(patient => patient && patient.id || null)
    );
    this.tags = this.patient.pipe(
      map(patient => patient && patient.tags || [])
    );
    this.categories = this.patient.pipe(
      map(patient => patient && patient.categories || [])
    );
  }

  ngOnInit() {
  }

}
