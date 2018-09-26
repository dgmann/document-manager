import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Record, RecordService } from "../core/store/index";
import { PatientService } from "./patient.service";
import { Patient } from "./store/patient.model";
import { Filter } from "./store/patient.reducer";
import { untilDestroyed } from "ngx-take-until-destroy";
import { Category, CategoryService } from "../core";


@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss']
})
export class PatientComponent implements OnInit, OnDestroy {
  public patientId: Observable<string>;
  public records: Observable<Record[]>;
  public patient: Observable<Patient>;
  public selectedRecord: Observable<Record>;
  public showDetails = false;
  public categories: Observable<{ [id: string]: Category }>;
  public selectedCategory: Observable<string>;

  constructor(private recordService: RecordService,
              private patientService: PatientService,
              private categoryService: CategoryService,
              private router: Router,
              private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params
      .pipe(untilDestroyed(this))
      .subscribe(params => this.patientService.selectPatient(params['id']));
    this.selectedCategory = this.patientService.getSelectedCategory();
    this.route.queryParamMap
      .pipe(untilDestroyed(this))
      .subscribe(params => this.patientService.selectCategory(params.get('category')));

    this.records = this.patientService.getFilteredPatientRecords();
    this.patient = this.patientService.getSelectedPatient();
    this.patientId = this.patient.pipe(
      map(patient => patient && patient.id || null)
    );
    this.categories = this.categoryService.getAsMap();

    this.selectedRecord = this.patientService.getSelectedRecord();
  }

  ngOnDestroy() {
  }

  onSelectRecord(id: string) {
    this.patientService.selectRecord(id);
    this.showDetails = true;
  }

  setFilter(filter: Filter) {
    this.patientService.setFilter(filter);
  }

  onSelectedCategoryChange(category: string) {
    this.router.navigate(["."], {relativeTo: this.route, queryParams: {category}});
  }
}
