import {AfterViewInit, Component, ElementRef, Inject, OnInit, ViewChild} from '@angular/core';
import {FormControl} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {take} from "rxjs/operators";
import {ReplaySubject} from "rxjs/ReplaySubject";


import {Record} from "../../store";
import {CategoryService} from "../category-service";
import {Patient, PatientService} from "../patient-service";
import {TagService} from "../tag-service";


@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss']
})
export class DocumentEditDialogComponent implements AfterViewInit, OnInit {
  @ViewChild('datepickertoogle', {read: ElementRef}) datepickerToggle;
  public record: Record;
  public tabIndex = new ReplaySubject<number>();
  public categoryFormControl: FormControl;

  constructor(public dialogRef: MatDialogRef<DocumentEditDialogComponent>,
              @Inject(MAT_DIALOG_DATA) record: Record,
              public patient: PatientService,
              public tagsService: TagService,
              public categoryService: CategoryService) {
    this.record = Object.assign({}, record);
    if (!this.record.date) {
      this.record.date = new Date();
    }
    this.record.tags = record.tags.slice();
    if (!this.record.patientId) {
      this.patient.getCurrent().pipe(take(1)).subscribe((patient: Patient) => this.record.patientId = this.record.patientId || patient.id);
    }
    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    }
    else {
      this.tabIndex.next(0);
    }
    this.categoryFormControl = new FormControl(this.record.category);
    this.categoryFormControl.valueChanges.subscribe(val => this.record.category = val);
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
