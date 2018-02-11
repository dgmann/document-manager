import {COMMA, ENTER} from "@angular/cdk/keycodes";
import {AfterViewInit, Component, ElementRef, Inject, OnInit, ViewChild} from '@angular/core';
import {MAT_DIALOG_DATA, MatAutocompleteSelectedEvent, MatChipInputEvent, MatDialogRef} from "@angular/material";
import {Observable} from "rxjs/Observable";
import {map, startWith, switchMap, take} from "rxjs/operators";
import {ReplaySubject} from "rxjs/ReplaySubject";
import {Subject} from "rxjs/Subject";
import {Patient, PatientService} from "../shared";
import {TagService} from "../shared/tag-service";


import {Record} from "../store";


@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss']
})
export class DocumentEditDialogComponent implements AfterViewInit, OnInit {
  @ViewChild('datepickertoogle', {read: ElementRef}) datepickerToggle;
  tagInput = new Subject<string>();

  availableTags: Observable<string[]>;
  filteredOptions: Observable<string[]>;

  selectable: boolean = true;
  removable: boolean = true;
  addOnBlur: boolean = false;

  separatorKeysCodes = [ENTER, COMMA];
  public record: Record;
  public tabIndex = new ReplaySubject<number>();

  constructor(public dialogRef: MatDialogRef<DocumentEditDialogComponent>,
              @Inject(MAT_DIALOG_DATA) record: Record,
              public patient: PatientService,
              public tagsService: TagService) {
    this.record = Object.assign({}, record);
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
    this.availableTags = tagsService.get();
  }

  ngOnInit() {
    this.filteredOptions = this.tagInput
      .pipe(
        startWith(''),
        switchMap((val: string) => this.filter(val))
      );
  }

  ngAfterViewInit(): void {
    this.datepickerToggle.nativeElement.querySelector('button').setAttribute('tabindex', '-1');
  }

  add(event: MatChipInputEvent, autocomplete): void {
    let input = event.input;
    let value = event.value;

    if (input) {
      this.resetInputValue(input);
    }
    if (autocomplete.isOpen) {
      return;
    }

    if ((value || '').trim()) {
      this.record.tags.push(value.trim());
    } else {
      this.dialogRef.close(this.record);
    }
  }

  addOption(event: MatAutocompleteSelectedEvent, tagInput): void {
    let value = event.option.value;

    if ((value || '').trim()) {
      this.record.tags.push(value.trim());
    }
    this.resetInputValue(tagInput);
  }

  remove(tag: string): void {
    const index = this.record.tags.indexOf(tag);
    if (index > -1) {
      this.record.tags.splice(index, 1);
    }
  }

  inputChanged(event) {
    let data = event.currentTarget.value;
    if (data) {
      this.tagInput.next(data);
    }
  }

  filter(val: string) {
    return this.availableTags.pipe(
      map(tags => tags.filter(option => option.toLowerCase().indexOf(val.toLowerCase()) === 0)),
      map(tags => tags.filter(tag => this.record.tags.findIndex(t => t.toLowerCase() == tag.toLowerCase()) < 0))
    );
  }

  resetInputValue(tagInput) {
    tagInput.value = "";
    this.tagInput.next("");
  }
}
