import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  Inject,
  OnDestroy,
  OnInit,
  ViewChild
} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import * as moment from 'moment';
import {Observable, of, ReplaySubject} from 'rxjs';
import {defaultIfEmpty, filter, map, mergeMap, startWith, take} from 'rxjs/operators';
import {Patient} from '@app/patient';


import {Record} from '@app/core/records';
import {TagService} from '@app/core/tags';
import {Category, CategoryService} from '@app/core/categories';
import {untilDestroyed} from 'ngx-take-until-destroy';
import {EditResult} from './edit-result.model';
import {ExternalApiService} from './external-api.service';

@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentEditDialogComponent implements AfterViewInit, OnInit, OnDestroy {
  @ViewChild('datepickertoogle', { read: ElementRef, static: true }) datepickerToggle;
  public record: Record;
  public tabIndex = new ReplaySubject<number>();
  public currentExternalPatient$: Observable<Patient>;
  public selectedPatient$: Observable<Patient>;

  public categories: Category[];
  public tags: Observable<string[]>;

  filteredCategories: Observable<Category[]>;

  editForm = new FormGroup({
    patient: new FormControl(''),
    date: new FormControl(moment()),
    category: new FormControl(null, Validators.required),
    tags: new FormControl()
  });

  constructor(public dialogRef: MatDialogRef<DocumentEditDialogComponent>,
              @Inject(MAT_DIALOG_DATA) record: Record,
              public patientService: ExternalApiService,
              public categoryService: CategoryService,
              public tagService: TagService) {
    this.record = Object.assign({}, record);
    this.record.tags = record.tags.slice();
  }

  ngOnInit() {
    this.categoryService.load();
    this.categoryService.categories.pipe(untilDestroyed(this)).subscribe(categories => this.categories = categories);

    this.filteredCategories = this.editForm.get('category').valueChanges
      .pipe(
        startWith(''),
        map(value => this._filter(value))
      );

    this.tagService.load();
    this.tags = this.tagService.tags;

    this.currentExternalPatient$ = this.patientService.getSelectedPatient();

    this.editForm.patchValue({
      date: this.record.date || moment(),
      category: this.record.category,
      tags: this.record.tags
    });

    if (this.record.patientId) {
      this.editForm.patchValue({
        patient: {id: this.record.patientId, firstName: '', lastName: ''} as Patient
      });
    }

    this.patientService.getPatientById(this.record.patientId).subscribe(patient => this.editForm.patchValue({
      patient,
    }));

    if (this.record.category) {
      this.categoryService.categoryMap.pipe(
        filter(categories => Object.entries(categories).length > 0),
        take(1)
      ).subscribe(categories => this.editForm.patchValue({
        category: categories[this.record.category],
      }));
    }

    this.selectedPatient$ = this.editForm.get('patient').valueChanges.pipe(
      startWith(this.editForm.get('patient').value),
      mergeMap(value => {
        if (!value) {
          return this.currentExternalPatient$;
        } else {
          return of(value);
        }
      })
    );

    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    } else {
      this.tabIndex.next(0);
    }
  }

  ngOnDestroy(): void {
  }

  ngAfterViewInit(): void {
    this.datepickerToggle.nativeElement.querySelector('button').setAttribute('tabindex', '-1');
  }

  onSubmit() {
    if (this.editForm.valid) {
      this.selectedPatient$.pipe(
        take(1),
        map(patient => ({
          id: this.record.id,
          change: {
            patientId: patient.id,
            date: this.editForm.get('date').value,
            tags: this.editForm.get('tags').value,
            category: this.editForm.get('category').value.id
          }
      }) as EditResult))
        .subscribe(changeSet => this.dialogRef.close(changeSet));
    }
  }

  displayFn(category: Category): string | undefined {
    return category ? category.name : undefined;
  }

  private _filter(value: any): Category[] {
    let filterValue = '';
    if (value.name) {
      filterValue = value.name.toLowerCase();
    } else {
      filterValue = value.toLowerCase();
    }

    if (!filterValue) {
      return this.categories;
    }

    return this.categories.filter(category => category.name.toLowerCase().includes(filterValue));
  }
}
