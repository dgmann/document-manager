import { Injectable } from "@angular/core";
import { ReplaySubject } from "rxjs/ReplaySubject";

@Injectable()
export class PatientService {
  private current: ReplaySubject<Patient>;

  constructor() {
    this.current = new ReplaySubject<Patient>();
    this.current.next({id: "3", name: "John Doe"});
  }

  getCurrent() {
    return this.current.asObservable();
  }
}

export class Patient {
  constructor(public id: string, public name: string) {
  }
}
