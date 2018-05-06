import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs";
import {Patient} from "..";
import {PatientService} from "../patient.service";

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent implements OnInit {
  patient: Observable<Patient>;

  constructor(private patientService: PatientService) {
  }

  ngOnInit() {
    this.patient = this.patientService.getSelectedPatient();
  }

}
