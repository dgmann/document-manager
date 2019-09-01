import {Injectable} from '@angular/core';
import {select, Store} from '@ngrx/store';
import {
  selectFilteredPatientRecords,
  selectPatientRecords,
  selectSelectedCategoryId,
  selectSelectedId,
  selectSelectedPatient,
  selectSelectedRecord,
  State
} from './reducers';
import {SelectCategory, SelectPatient, SelectRecord, SetFilter} from './store/patient.actions';
import {Filter} from './store/patient.reducer';
import {Observable} from 'rxjs';
import {Patient} from './store/patient.model';
import {Record} from '../core/store/record';

@Injectable()
export class PatientService {
  public selectedPatient$: Observable<Patient>;
  public selectedPatientRecords$: Observable<Record[]>;
  public selectedRecord$: Observable<Record>;
  public selectedCategory$: Observable<string>;
  public filteredPatientRecord$: Observable<Record[]>;
  public selectedId$: Observable<string>;

  constructor(private store: Store<State>) {
    this.selectedPatient$ = this.store.pipe(select(selectSelectedPatient));
    this.selectedPatientRecords$ = this.store.pipe(select(selectPatientRecords));
    this.selectedRecord$ = this.store.pipe(select(selectSelectedRecord));
    this.selectedCategory$ = this.store.pipe(select(selectSelectedCategoryId));
    this.filteredPatientRecord$ = this.store.pipe(select(selectFilteredPatientRecords));
    this.selectedId$ = this.store.pipe(select(selectSelectedId));
  }

  public selectPatient(id: string) {
    this.store.dispatch(new SelectPatient({id}));
  }

  public setFilter(filter: Filter) {
    this.store.dispatch(new SetFilter(filter));
  }

  public selectRecord(id: string) {
    this.store.dispatch(new SelectRecord({id}));
  }

  public selectCategory(id: string) {
    this.store.dispatch(new SelectCategory({id}));
  }
}
