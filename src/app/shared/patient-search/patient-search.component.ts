import {Component, OnInit} from '@angular/core';
import {FormControl} from "@angular/forms";
import {Observable} from "rxjs/Observable";
import {debounceTime, filter, switchMap} from "rxjs/operators";
import {Patient, PatientService} from "../patient-service";

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss']
})
export class PatientSearchComponent implements OnInit {
  public searchResults: Observable<Patient[]>;
  public searchInput = new FormControl();

  constructor(private patientService: PatientService) {
  }

  ngOnInit() {
    this.searchResults = this.searchInput.valueChanges
      .pipe(
        debounceTime(500),
        filter(query => !!query && query.length > 0),
        switchMap(query => this.patientService.find(query))
      );
  }

  displayFn(patient: Patient): string | undefined {
    return patient ? patient.name : undefined;
  }

}
