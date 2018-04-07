import {AfterViewInit, Component, ElementRef, Inject, OnInit, ViewChild} from '@angular/core';
import {FormControl} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {take} from "rxjs/operators";
import {ReplaySubject} from "rxjs/ReplaySubject";
import {Patient} from "../../patient";


import {Record} from "../../store";
import {CategoryService} from "../category-service";
import {ExternalApiService} from "../external-api.service";
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
              public patient: ExternalApiService,
              public tagsService: TagService,
              public categoryService: CategoryService) {
    this.record = Object.assign({}, record);
    if (!this.record.date) {
      this.record.date = new Date();
    }
    this.record.tags = record.tags.slice();
    if (!this.record.patientId) {
      this.patient.getSelectedPatient().pipe(take(1)).subscribe((patient: Patient) => this.record.patientId = this.record.patientId || patient.id);
    }
    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    }
    else {
      this.tabIndex.next(0);
    }
    this.categoryFormControl = new FormControl(this.record.categoryId);
    this.categoryFormControl.valueChanges.subscribe(val => this.record.categoryId = val.id);
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
    return this.record.date && this.record.patientId && this.record.categoryId;
  }
}
