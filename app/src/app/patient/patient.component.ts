import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { Observable } from "rxjs";
import { filter, switchMap } from "rxjs/operators";
import { Record, RecordService } from "../core/store";
import { PatientService } from "./patient.service";
import { Patient } from "./store/patient.model";
import { Filter } from "./store/patient.reducer";
import { untilDestroyed } from "ngx-take-until-destroy";
import { Category, CategoryService, TagService } from "../core";
import { EditResult } from "../shared";


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

  public availableCategories: Observable<Category[]>;
  public availableTags: Observable<string[]>;

  constructor(private recordService: RecordService,
              private patientService: PatientService,
              private categoryService: CategoryService,
              private tagsService: TagService,
              private router: Router,
              private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params
      .pipe(untilDestroyed(this))
      .subscribe(params => this.patientService.selectPatient(params['id']));
    this.selectedCategory = this.patientService.selectedCategory$;
    this.route.queryParamMap
      .pipe(untilDestroyed(this))
      .subscribe(params => this.patientService.selectCategory(params.get('category')));

    this.records = this.patientService.filteredPatientRecord$;
    this.patient = this.patientService.selectedPatient$;
    this.patientId = this.patientService.selectedId$;
    this.categoryService.load();
    this.categories = this.categoryService.categoryMap;

    this.availableCategories = this.patientId.pipe(filter(p => !!p), switchMap(id => this.categoryService.getByPatientId(id)));
    this.availableTags = this.patientId.pipe(filter(p => !!p), switchMap(id => this.tagsService.getByPatientId(id)));

    this.selectedRecord = this.patientService.selectedRecord$;
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

  onUpdateRecord(data: EditResult) {
    this.recordService.update(data.id, data.change);
  }
}
