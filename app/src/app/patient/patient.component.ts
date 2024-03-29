import { Component, DestroyRef, inject, OnDestroy, OnInit } from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Observable} from 'rxjs';
import {filter, switchMap} from 'rxjs/operators';
import {Record, RecordService} from '@app/core/records';
import {PatientService} from './patient.service';
import {Patient} from './store/patient.model';
import {Filter} from './store/patient.reducer';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import {Category, CategoryService} from '@app/core/categories';
import {EditResult} from '../shared';
import {animate, AnimationEvent, style, transition, trigger} from '@angular/animations';
import {TagService} from '@app/core/tags';


@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss'],
  animations: [
    trigger('panelInOut', [
      transition('void => *', [
        style({transform: 'translateX(100%)'}),
        animate(200)
      ]),
      transition('* => void', [
        animate(200, style({transform: 'translateX(100%)'}))
      ])
    ])
  ]
})
export class PatientComponent implements OnInit, OnDestroy {
  public patientId$: Observable<string>;
  public records$: Observable<Record[]>;
  public patient$: Observable<Patient>;
  public selectedRecord$: Observable<Record>;
  public categories$: Observable<{ [id: string]: Category }>;

  public availableCategories$: Observable<Category[]>;
  public availableTags$: Observable<string[]>;

  destroyRef = inject(DestroyRef);

  constructor(private recordService: RecordService,
              private patientService: PatientService,
              private categoryService: CategoryService,
              private tagsService: TagService,
              private router: Router,
              private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.params
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe(params => this.patientService.selectPatient(params.id));

    this.records$ = this.patientService.filteredPatientRecord$;
    this.patient$ = this.patientService.selectedPatient$;
    this.patientId$ = this.patientService.selectedId$;
    this.categoryService.load();
    this.categories$ = this.categoryService.categoryMap;

    this.availableCategories$ = this.patientId$.pipe(filter(p => !!p), switchMap(id => this.categoryService.getByPatientId(id)));
    this.availableTags$ = this.patientId$.pipe(filter(p => !!p), switchMap(id => this.tagsService.getByPatientId(id)));

    this.selectedRecord$ = this.patientService.selectedRecord$;
  }

  ngOnDestroy() {
  }

  onSelectRecord(id: string) {
    this.patientService.selectRecord(id);
  }

  setFilter(set: Filter) {
    this.patientService.setFilter(set);
  }

  onUpdateRecord(data: EditResult) {
    this.recordService.update(data.id, data.change);
  }

  onDeleteRecord(record: Record) {
    this.recordService.delete(record.id);
  }

  onOpenInEditor(record: Record) {
    this.router.navigate(['/editor', record.id]);
  }

  onDuplicateRecord(record: Record) {
    this.recordService.duplicate(record.id);
  }

  onAnimationEvent(event: AnimationEvent) {
    event.element.focus();
  }

  onDetailsPanelClose() {
    this.patientService.selectRecord(null);
  }
}
