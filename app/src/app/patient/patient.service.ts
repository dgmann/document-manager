import {Injectable} from "@angular/core";
import {select, Store} from "@ngrx/store";
import {selectPatientRecords, selectSelectedPatient, State} from "./reducers";
import {SelectPatient} from "./store/patient.actions";

@Injectable()
export class PatientService {
  constructor(private store: Store<State>) {
  }

  public selectPatient(id: string) {
    this.store.dispatch(new SelectPatient({id: id}));
  }

  public getSelectedPatient() {
    return this.store.pipe(select(selectSelectedPatient));
  }

  public getPatientRecords() {
    return this.store.pipe(select(selectPatientRecords));
  }
}
