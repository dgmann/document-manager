import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {ReplaySubject} from "rxjs/ReplaySubject";
import {environment} from "../../environments/environment";
import {Category} from "./category-service";

@Injectable()
export class PatientService {
  private current: ReplaySubject<Patient>;

  constructor(private http: HttpClient) {
    this.current = new ReplaySubject<Patient>();
    this.current.next({id: "3", firstName: "John", lastName: "Doe", birthDate: new Date()});
  }

  getCurrent() {
    return this.current.asObservable();
  }

  find(query: string) {
    return this.http.get<Patient[]>(`${environment.api}/patients?name=${query}`);
  }
}

export interface Patient {
  id: string;
  firstName: string;
  lastName: string;
  birthDate: Date;
  tags?: string[];
  categories?: Category[];
}
