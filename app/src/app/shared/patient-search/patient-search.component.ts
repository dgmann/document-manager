import { HttpClient } from "@angular/common/http";
import { ChangeDetectionStrategy, Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormControl } from "@angular/forms";
import { MatAutocompleteSelectedEvent } from "@angular/material";
import { Observable } from "rxjs";
import { debounceTime, filter, switchMap, map } from "rxjs/operators";
import { environment } from "../../../environments/environment";
import { Patient } from "../../patient";

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
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
        map(query => this.parseQuery(query)),
        switchMap(query => this.http.get<Patient[]>(`${environment.api}/patients`, {params: {...query}}))
      );
  }
      
  parseQuery(query: string) {
    const parts = query.split(",");
    const result {
      lastname: parts[0],
      fistname: parts[1]
    };
    if (!result.firstname) {
      delete result.firstname;
    }
    return result;
  }

  displayFn(patient: Patient): string | undefined {
    return patient ? patient.lastName + ', ' + patient.firstName : undefined;
  }

  onSelectPatient(event: MatAutocompleteSelectedEvent) {
    this.selectPatient.emit(event.option.value);
    this.searchInput.reset();
  }

}
