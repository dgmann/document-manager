import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {ReplaySubject} from "rxjs/ReplaySubject";
import {environment} from "../../environments/environment";

@Injectable()
export class PatientService {
  private current: ReplaySubject<Patient>;

  constructor(private http: HttpClient) {
    this.current = new ReplaySubject<Patient>();
    this.current.next({id: "3", name: "John Doe"});
  }

  getCurrent() {
    return this.current.asObservable();
  }

  find(query: string) {
    return this.http.get<Patient[]>(`${environment.api}/patients?query=${query}`);
  }
}

export class Patient {
  constructor(public id: string, public name: string) {
  }
}
