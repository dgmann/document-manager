import { COMMA, ENTER } from "@angular/cdk/keycodes";
import { ChangeDetectionStrategy, Component, forwardRef, Input, OnInit, ViewChild } from '@angular/core';
import { ControlValueAccessor, FormControl, NG_VALUE_ACCESSOR } from "@angular/forms";
import { MatAutocomplete, MatChipInputEvent, MatInput } from "@angular/material";
import { difference } from 'lodash-es';
import { Observable } from "rxjs";
import { map, startWith } from "rxjs/operators";

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
export class AutocompleteChipsComponent implements OnInit, ControlValueAccessor {
  @Input('options') options: string[];
  @ViewChild('chipInput') chipInput: MatInput;
  @ViewChild('auto') autoComplete: MatAutocomplete;

  values: string[] = [];
  formControl = new FormControl();

  filteredOptions: Observable<string[]>;
  separatorKeysCodes = [ENTER, COMMA];

  propagateChange = (_: any) => {
  };

  ngOnInit() {
    const filtered = this.formControl.valueChanges
      .pipe(
        startWith(''),
        map((val => difference(this.filter(val, this.options), this.values)))
      );
    this.filteredOptions = filtered;
  }

  filter(val: string, options: string[]): string[] {
    if (!val) {
      return options;
    }
    const result = options.filter(option =>
      option.toLowerCase().indexOf(val.toLowerCase()) === 0);
    return result;
  }

  addValue(value: string) {
    if ((value || '').trim()) {
      this.values.push(value.trim());
      this.propagateChange(this.values);
    }

    this.reset();
  }

  addChip(event: MatChipInputEvent): void {
    let value = this.formControl.value;
    this.addValue(value);
  }

  remove(chip: any): void {
    let index = this.values.indexOf(chip);

    if (index >= 0) {
      this.values.splice(index, 1);
      this.propagateChange(this.values);
    }
  }

  reset() {
    this.chipInput['nativeElement'].value = '';
    this.formControl.reset();
  }

  onSubmit(event) {
    if (this.chipInput['nativeElement'].value == '' && !this.autoComplete.isOpen) {
      event.preventDefault();
    }
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
