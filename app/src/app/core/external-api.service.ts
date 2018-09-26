import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { ReplaySubject } from "rxjs";
import { Patient } from "../patient";

@Injectable({
  providedIn: "root"
})
export class ExternalApiService {
  private current: ReplaySubject<Patient>;

  constructor(private http: HttpClient) {
    this.current = new ReplaySubject<Patient>();
    this.current.next({id: "3", firstName: "John", lastName: "Doe", birthDate: new Date()});
  }

  getSelectedPatient() {
    return this.http.get<Patient>("http://localhost:3000/patient");
  }
}
