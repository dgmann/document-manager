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
import { FormControl, FormGroup } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material";
import * as moment from "moment";
import { Observable, ReplaySubject } from "rxjs";
import { filter, map, startWith, take } from "rxjs/operators";
import { Patient } from "../../patient";


import { Record } from "../../core/store";
import { Category, CategoryService, ExternalApiService, TagService } from "../../core";
import { untilDestroyed } from "ngx-take-until-destroy";
import { EditResult } from "./edit-result.model";

@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentEditDialogComponent implements AfterViewInit, OnInit, OnDestroy {
  @ViewChild('datepickertoogle', {read: ElementRef}) datepickerToggle;
  public record: Record;
  public tabIndex = new ReplaySubject<number>();
  public patient: Patient;

  public categories: Category[];
  public tags: Observable<string[]>;

  filteredCategories: Observable<Category[]>;

  editForm = new FormGroup({
    patientId: new FormControl(''),
    date: new FormControl(moment()),
    category: new FormControl(),
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

    this.patientService.getSelectedPatient().subscribe(p => {
      this.patient = p;
      if (!this.editForm.get('patientId').value) {
        this.editForm.patchValue({
          patientId: this.record.patientId,
        })
      }
    });

    this.editForm.patchValue({
      patientId: this.record.patientId,
      date: this.record.date || moment(),
      category: this.record.category,
      tags: this.record.tags
    });

    if (this.record.category) {
      this.categoryService.categoryMap.pipe(
        filter(categories => Object.entries(categories).length > 0),
        take(1)
      ).subscribe(categories => this.editForm.patchValue({
        category: categories[this.record.category],
      }));
    }

    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    }
    else {
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
      const changeSet: EditResult = {
        id: this.record.id, change: {
          patientId: this.editForm.get('patientId').value,
          date: this.editForm.get('date').value,
          tags: this.editForm.get('tags').value,
          category: this.editForm.get('category').value.id
        }
      };
      this.dialogRef.close(changeSet);
    }
  }

  displayFn(category: Category): string | undefined {
    return category ? category.name : undefined;
  }

  private _filter(value: any): Category[] {
    let filterValue = "";
    if (value.name) {
      filterValue = value.name.toLowerCase();
    } else {
      filterValue = value.toLowerCase();
    }

    return this.categories.filter(category => category.name.toLowerCase().includes(filterValue));
  }
}
