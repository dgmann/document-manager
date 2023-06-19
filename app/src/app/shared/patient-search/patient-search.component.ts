import {FocusMonitor} from '@angular/cdk/a11y';
import {coerceBooleanProperty} from '@angular/cdk/coercion';
import {HttpClient} from '@angular/common/http';
import {
  ChangeDetectionStrategy,
  Component, ElementRef,
  EventEmitter,
  HostBinding, HostListener,
  Input, OnDestroy,
  OnInit,
  Optional,
  Output,
  Self
} from '@angular/core';
import {ControlValueAccessor, UntypedFormControl, NgControl} from '@angular/forms';
import {MatAutocompleteSelectedEvent} from '@angular/material/autocomplete';
import {MatFormFieldControl} from '@angular/material/form-field';
import {ExternalApiService} from '@app/shared/document-edit-dialog/external-api.service';
import {isEqual} from 'lodash-es';
import { Observable, EMPTY, Subject, of } from 'rxjs';
import {catchError, debounceTime, filter, map, startWith, switchMap} from 'rxjs/operators';
import {Patient} from '@app/patient';

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [{provide: MatFormFieldControl, useExisting: PatientSearchComponent}],
})
export class PatientSearchComponent implements OnInit, OnDestroy, ControlValueAccessor, MatFormFieldControl<Patient | string> {
  static nextId = 0;
  @HostBinding() id = `app-patient-search-${PatientSearchComponent.nextId++}`;

  @Output() selectPatient = new EventEmitter<Patient>();
  @Input() set value(value: Patient | string) {
    const setPatient = (v: Patient) => {
      this._value = v;
      this.onChange(v);
      this.searchInput.setValue(v);
      this.stateChanges.next();
    };

    // Value could also be just a patient id in some cases
    const patientId = parseInt(value as string, 10);
    if (patientId) {
      this.patientService.getPatientById(patientId + '').subscribe(patient => setPatient(patient));
    } else {
      setPatient(value as Patient);
    }
  }
  get value() {
    return this._value;
  }
  // tslint:disable-next-line:variable-name
  private _value: Patient;

  @Input() clearOnSubmit = false;

  @Input()
  get placeholder() {
    return this._placeholder;
  }
  set placeholder(plh) {
    this._placeholder = plh;
    this.stateChanges.next();
  }
  // tslint:disable-next-line:variable-name
  private _placeholder: string;

  autofilled = false;
  controlType = 'app-patient-search';
  get empty() {
    return !this.searchInput.value;
  }
  readonly errorState = false;

  focused = false;

  @Input()
  get required() {
    return this._required;
  }
  set required(req) {
    this._required = coerceBooleanProperty(req);
    this.stateChanges.next();
  }
  // tslint:disable-next-line:variable-name
  private _required = false;

  @Input()
  get disabled() {
    return this._disabled;
  }
  set disabled(dis) {
    this._disabled = coerceBooleanProperty(dis);
    this.stateChanges.next();
  }
  // tslint:disable-next-line:variable-name
  private _disabled = false;

  @HostBinding('class.floating')
  get shouldLabelFloat() {
    return this.focused || !this.empty;
  }
  stateChanges = new Subject<void>();
  userAriaDescribedBy: string;

  searchResults: Observable<Patient[]>;
  searchInput = new UntypedFormControl();

  onChange: (patient: Patient) => void = (p) => {};
  onTouched = () => {};

  constructor(private patientService: ExternalApiService, private fm: FocusMonitor, @Optional() @Self() public ngControl: NgControl, private elRef: ElementRef) {
    // Setting the value accessor directly (instead of using
    // the providers) to avoid running into a circular import.
    if (this.ngControl != null) { this.ngControl.valueAccessor = this; }

    fm.monitor(elRef.nativeElement, true).subscribe(origin => {
      this.focused = !!origin;
      this.stateChanges.next();
    });
  }

  ngOnInit() {
    this.searchResults = this.searchInput.valueChanges
      .pipe(
        startWith(this.searchInput.value),
        debounceTime(500),
        filter(query => !!query && query.length > 0),
        switchMap(query => {
          const patientId = parseInt(query, 10);
          if (patientId) {
            return this.patientService.getPatientById(query).pipe(
              map(patient => [patient]),
              catchError(err => of<Patient[]>([{id: query, firstName: 'unknown', lastName: 'unknown', birthDate: new Date()}]))
            );
          } else {
            const patientQuery = this.parseQuery(query);
            return this.patientService.find(patientQuery);
          }
        })
      );

    // Unset value as it may not match the currently selected value anymore
    this.searchInput.valueChanges.subscribe(current => {
      if (this.value && !isEqual(this.value, current)) {
        this.value = null;
      }
    });
  }

  parseQuery(query: string) {
    const parts = query.split(',');
    const result = {
      lastname: parts[0] && parts[0].trim() || undefined,
      firstname: parts[1] && parts[1].trim() || undefined
    };
    if (!result.firstname) {
      delete result.firstname;
    }
    return result;
  }

  displayFn(patient: Patient): string | undefined {
    return patient ? patient.lastName + ', ' + patient.firstName : undefined;
  }

  onSelectPatient(event: MatAutocompleteSelectedEvent) {
    this.value = event.option.value;
    if (this.clearOnSubmit) {
      this.searchInput.reset();
    }
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    if (isDisabled) {
      this.searchInput.disable();
    } else {
      this.searchInput.enable();
    }
  }

  writeValue(value: any): void {
    this.value = value;
  }

  onContainerClick(event: MouseEvent) {
    if ((event.target as Element).tagName.toLowerCase() !== 'input') {
      this.elRef.nativeElement.querySelector('input').focus();
    }
  }

  ngOnDestroy() {
    this.stateChanges.complete();
    this.fm.stopMonitoring(this.elRef.nativeElement);
  }

  @HostBinding('attr.aria-describedby') describedBy = '';

  setDescribedByIds(ids: string[]) {
    this.describedBy = ids.join(' ');
  }

  @HostListener('focus', ['$event'])
  onFocus(e) {
    this.elRef.nativeElement.querySelector('input').focus();
  }

  onEmptied(e) {
    console.log(e);
  }
}
