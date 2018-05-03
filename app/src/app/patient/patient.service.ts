import { Injectable } from "@angular/core";
import { select, Store } from "@ngrx/store";
import {
  selectFilteredPatientRecords,
  selectPatientRecords,
  selectSelectedPatient,
  selectSelectedRecord,
  State
} from "./reducers";
import { SelectPatient, SelectRecord, SetFilter } from "./store/patient.actions";

@Injectable()
export class PatientService {
  constructor(private store: Store<State>) {
  }

  public selectPatient(id: string) {
    this.store.dispatch(new SelectPatient({id: id}));
  }

  public setFilter(categoryIds: string[], tags: string[]) {
    this.store.dispatch(new SetFilter({categoryIds: categoryIds, tags: tags}));
  }

  public getSelectedPatient() {
    return this.store.pipe(select(selectSelectedPatient));
  }

  public getPatientRecords() {
    return this.store.pipe(select(selectPatientRecords));
  }

  public getFilteredPatientRecords() {
    return this.store.pipe(select(selectFilteredPatientRecords));
  }

  public selectRecord(id: string) {
    this.store.dispatch(new SelectRecord({id: id}));
  }

  public getSelectedRecord() {
    return this.store.pipe(select(selectSelectedRecord));
  }
}
