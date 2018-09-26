import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  Inject,
  OnInit,
  ViewChild
} from '@angular/core';
import { FormControl } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material";
import * as moment from "moment";
import { Observable, ReplaySubject } from "rxjs";
import { shareReplay } from "rxjs/operators";
import { Patient } from "../../patient";


import { Record } from "../../core/store";
import { Category, CategoryService, ExternalApiService } from "../../core";

@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DocumentEditDialogComponent implements AfterViewInit, OnInit {
  @ViewChild('datepickertoogle', {read: ElementRef}) datepickerToggle;
  public record: Record;
  public tabIndex = new ReplaySubject<number>();
  public categoryFormControl: FormControl;
  public patient: Patient;
  categories: Observable<Category[]>;

  constructor(public dialogRef: MatDialogRef<DocumentEditDialogComponent>,
              @Inject(MAT_DIALOG_DATA) record: Record,
              public patientService: ExternalApiService,
              public categoryService: CategoryService) {
    let patientRequest = this.patientService.getSelectedPatient().pipe(shareReplay());
    this.categories = categoryService.get();
    patientRequest.subscribe(p => this.patient = p);
    this.record = Object.assign({}, record);
    if (!this.record.date) {
      this.record.date = moment();
    }
    this.record.tags = record.tags.slice();
    if (!this.record.patientId) {
      patientRequest.subscribe(p => this.record.patientId = p.id);
    }
    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    }
    else {
      this.tabIndex.next(0);
    }
    this.categoryFormControl = new FormControl({name: "Test", id: this.record.category});
    this.categoryFormControl.valueChanges.subscribe(val => this.record.category = val.id);
  }

  ngOnInit() {

  }

  ngAfterViewInit(): void {
    this.datepickerToggle.nativeElement.querySelector('button').setAttribute('tabindex', '-1');
  }

  submit() {
    if (this.valid()) {
      this.dialogRef.close(this.record);
    }
  }

  valid() {
    return this.record.date && this.record.patientId && this.record.category;
  }
}
