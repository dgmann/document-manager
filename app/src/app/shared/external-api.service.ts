import {Injectable} from "@angular/core";
import {ReplaySubject} from "rxjs";
import {Patient} from "../patient/store/patient.model";

@Injectable()
export class ExternalApiService {
  private current: ReplaySubject<Patient>;

  constructor() {
    this.current = new ReplaySubject<Patient>();
    this.current.next({id: "3", firstName: "John", lastName: "Doe", birthDate: new Date()});
  }

  getSelectedPatient() {
    return this.current.asObservable();
  }
}
