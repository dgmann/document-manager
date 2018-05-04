import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { Observable } from "rxjs/Observable";
import { map } from "rxjs/operators";
import { Record, RecordService } from "../store";
import { PatientService } from "./patient.service";
import { Patient } from "./store/patient.model";
import { Filter } from "./store/patient.reducer";

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss']
})
export class PatientComponent implements OnInit {
  public patientId: Observable<string>;
  public records: Observable<Record[]>;
  public patient: Observable<Patient>;
  public selectedRecord: Observable<Record>;
  public showDetails = false;

  constructor(private recordService: RecordService,
              private patientService: PatientService,
              private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params.subscribe(params => this.patientService.selectPatient(params['id']));
    this.records = this.patientService.getFilteredPatientRecords();
    this.patient = this.patientService.getSelectedPatient();
    this.patientId = this.patient.pipe(
      map(patient => patient && patient.id || null)
    );


    this.selectedRecord = this.patientService.getSelectedRecord();
  }

  onSelectRecord(id: string) {
    this.patientService.selectRecord(id);
    this.showDetails = true;
  }

  setFilter(filter: Filter) {
    this.patientService.setFilter(filter);
  }
}
