import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {ConfigService} from '@app/core/config';
import {EMPTY, Observable, of, ReplaySubject} from 'rxjs';
import {Patient} from '@app/patient';

@Injectable({
  providedIn: 'root'
})
export class ExternalApiService {

  constructor(private http: HttpClient, private config: ConfigService) {
  }

  getSelectedPatient() {
    return of({id: '1', firstName: 'Jon', lastName: 'Doe', birthDate: new Date('01.01.2000')} as Patient);
    //return this.http.get<Patient>('http://localhost:3000/patient');
  }

  getPatientById(id: string): Observable<Patient> {
    if (!id) {
      return EMPTY;
    }

    return this.http.get<Patient>(`${this.config.getApiUrl()}/patients/${id}`);
  }

  find(query: {lastname: string, firstname: string}): Observable<Patient[]> {
    return this.http.get<Patient[]>(`${this.config.getApiUrl()}/patients`, {params: {...query}});
  }
}
