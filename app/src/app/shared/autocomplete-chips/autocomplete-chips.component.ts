import {COMMA, ENTER} from '@angular/cdk/keycodes';
import {ChangeDetectionStrategy, Component, ElementRef, forwardRef, Input, ViewChild} from '@angular/core';
import {ControlValueAccessor, UntypedFormControl, NG_VALUE_ACCESSOR} from '@angular/forms';
import {MatAutocomplete, MatAutocompleteSelectedEvent} from '@angular/material/autocomplete';
import {MatChipInputEvent, MatChipListbox} from '@angular/material/chips';
import {difference} from 'lodash-es';
import {Observable} from 'rxjs';
import {map, startWith} from 'rxjs/operators';

@Component({
  selector: 'app-autocomplete-chips',
  templateUrl: './autocomplete-chips.component.html',
  styleUrls: ['./autocomplete-chips.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => AutocompleteChipsComponent),
      multi: true
    }
  ]
})
export class AutocompleteChipsComponent implements ControlValueAccessor {
  @Input() options: string[];
  @ViewChild('chipInput', {static: true}) chipInput: ElementRef<HTMLInputElement>;
  @ViewChild('chips', {static: true}) chips: MatChipListbox;
  @ViewChild('auto', { static: true }) autoComplete: MatAutocomplete;

  values: string[] = [];
  formControl = new UntypedFormControl();

  filteredOptions: Observable<string[]>;
  separatorKeysCodes = [ENTER, COMMA];

  propagateChange = (_: any) => {
  }

  constructor() {
    this.filteredOptions = this.formControl.valueChanges.pipe(
      startWith(null),
      map((v: string | null) => v ? this._filter(v) : difference(this.options.slice(), this.values)));
  }

  add(event: MatChipInputEvent): void {
    if (!this.autoComplete.isOpen) {
      const input = event.input;
      const value = event.value;

      if ((value || '').trim()) {
        this.values.push(value.trim());
        this.propagateChange(this.values);
      }

      // Reset the input value
      if (input) {
        input.value = '';
      }

      this.formControl.setValue(null);
    }
  }

  remove(fruit: string): void {
    const index = this.values.indexOf(fruit);

    if (index >= 0) {
      this.values.splice(index, 1);
      this.propagateChange(this.values);
    }
  }

  selected(event: MatAutocompleteSelectedEvent): void {
    this.values.push(event.option.viewValue);
    this.chipInput.nativeElement.value = '';
    this.formControl.setValue(null);
    this.propagateChange(this.values);
  }

  private _filter(value: string): string[] {
    const filterValue = value.toLowerCase();

    return this.options.filter(v => v.toLowerCase().indexOf(filterValue) === 0 && !this.values.includes(v));
  }

  registerOnChange(fn: any): void {
    this.propagateChange = fn;
  }

  registerOnTouched(fn: any): void {
  }

  setDisabledState(isDisabled: boolean): void {
  }

  writeValue(obj: any): void {
    if (obj !== undefined && obj !== null) {
      this.values = obj;
    }
  }
}
