import {HttpClient} from "@angular/common/http";
import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {FormControl} from "@angular/forms";
import {MatAutocompleteSelectedEvent} from "@angular/material";
import {Observable} from "rxjs/Observable";
import {debounceTime, filter, switchMap} from "rxjs/operators";
import {environment} from "../../../environments/environment";
import {Patient} from "../../patient";

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss']
})
export class PatientSearchComponent implements OnInit {
  @Output('selectPatient') selectPatient = new EventEmitter<Patient>();
  public searchResults: Observable<Patient[]>;
  public searchInput = new FormControl();

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    this.searchResults = this.searchInput.valueChanges
      .pipe(
        debounceTime(500),
        filter(query => !!query && query.length > 0),
        switchMap(query => this.http.get<Patient[]>(`${environment.api}/patients?name=${query}`))
      );
  }

  displayFn(patient: Patient): string | undefined {
    return patient ? patient.lastName + ', ' + patient.firstName : undefined;
  }

  onSelectPatient(event: MatAutocompleteSelectedEvent) {
    this.selectPatient.emit(event.option.value);
  }

}
