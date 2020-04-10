import {HttpClient} from '@angular/common/http';
import {ChangeDetectionStrategy, Component, EventEmitter, OnInit, Output} from '@angular/core';
import {FormControl} from '@angular/forms';
import {MatAutocompleteSelectedEvent} from '@angular/material/autocomplete';
import {Observable} from 'rxjs';
import {debounceTime, filter, map, switchMap} from 'rxjs/operators';
import {Patient} from '@app/patient';
import {ConfigService} from '@app/core/config';

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class PatientSearchComponent implements OnInit {
  @Output() selectPatient = new EventEmitter<Patient>();
  public searchResults: Observable<Patient[]>;
  public searchInput = new FormControl();

  constructor(private http: HttpClient, private config: ConfigService) {
  }

  ngOnInit() {
    this.searchResults = this.searchInput.valueChanges
      .pipe(
        debounceTime(500),
        filter(query => !!query && query.length > 0),
        switchMap(query => {
          const patientId = parseInt(query, 10);
          if (patientId) {
            return this.http.get<Patient>(`${this.config.getApiUrl()}/patients/${patientId}`).pipe(
              map(patient => [patient])
            );
          } else {
            const patientQuery = this.parseQuery(query);
            return this.http.get<Patient[]>(`${this.config.getApiUrl()}/patients`, {params: {...patientQuery}});
          }
        })
      );
  }

  parseQuery(query: string) {
    const parts = query.split(',');
    const result = {
      lastname: parts[0] && parts[0].trim() || undefined,
      firstname: parts[1] && parts[1].trim() || undefined
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
