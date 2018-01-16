import { AfterViewInit, Component, ElementRef, Inject, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA, MatChipInputEvent, MatDialogRef } from "@angular/material";
import { COMMA, ENTER } from "@angular/cdk/keycodes";


import { Record } from "../api";
import { Patient, PatientService } from "../shared";
import { ReplaySubject } from "rxjs/ReplaySubject";


@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss']
})
export class DocumentEditDialogComponent implements AfterViewInit {
  @ViewChild('datepickertoogle', {read: ElementRef}) datepickerToggle;

  selectable: boolean = true;
  removable: boolean = true;
  addOnBlur: boolean = false;

  separatorKeysCodes = [ENTER, COMMA];
  public record: Record;
  public tabIndex = new ReplaySubject<number>();

  constructor(public dialogRef: MatDialogRef<DocumentEditDialogComponent>, @Inject(MAT_DIALOG_DATA) record: Record, public patient: PatientService) {
    this.record = Object.assign({}, record);
    this.record.tags = record.tags.slice();
    if (!this.record.patientId) {
      this.patient.getCurrent().take(1).subscribe((patient: Patient) => this.record.patientId = this.record.patientId || patient.id);
    }
    if (this.record.patientId && this.record.date) {
      this.tabIndex.next(-1);
    }
    else {
      this.tabIndex.next(0);
    }
  }

  ngAfterViewInit(): void {
    this.datepickerToggle.nativeElement.querySelector('button').setAttribute('tabindex', '-1');
  }

  add(event: MatChipInputEvent): void {
    let input = event.input;
    let value = event.value;

    if ((value || '').trim()) {
      this.record.tags.push(value.trim());
    } else {
      this.dialogRef.close(this.record);
    }
    if (input) {
      input.value = '';
    }
  }

  remove(tag: string): void {
    const index = this.record.tags.indexOf(tag);
    if (index > -1) {
      this.record.tags.splice(index, 1);
    }
  }
}
