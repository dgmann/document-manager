import {COMMA, ENTER} from '@angular/cdk/keycodes';
import {ChangeDetectionStrategy, Component, forwardRef, Input, OnInit, ViewChild} from '@angular/core';
import {ControlValueAccessor, FormControl, NG_VALUE_ACCESSOR} from '@angular/forms';
import {MatAutocomplete} from '@angular/material/autocomplete';
import {MatChipInputEvent} from '@angular/material/chips';
import {MatInput} from '@angular/material/input';
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
export class AutocompleteChipsComponent implements OnInit, ControlValueAccessor {
  @Input('options') options: string[];
  @ViewChild('chipInput', { static: true }) chipInput: MatInput;
  @ViewChild('auto', { static: true }) autoComplete: MatAutocomplete;

  values: string[] = [];
  formControl = new FormControl();

  filteredOptions: Observable<string[]>;
  separatorKeysCodes = [ENTER, COMMA];

  propagateChange = (_: any) => {
  };

  ngOnInit() {
    this.filteredOptions = this.formControl.valueChanges
      .pipe(
        startWith(''),
        map((val => difference(this.filter(val, this.options), this.values)))
      );
  }

  filter(val: string, options: string[]): string[] {
    if (!val) {
      return options;
    }
    return options.filter(option =>
      option.toLowerCase().indexOf(val.toLowerCase()) === 0);
  }

  addValue(value: string) {
    if ((value || '').trim()) {
      this.values.push(value.trim());
      this.propagateChange(this.values);
    }

    this.reset();
  }

  addChip(event: MatChipInputEvent): void {
    const value = this.formControl.value;
    this.addValue(value);
  }

  remove(chip: any): void {
    const index = this.values.indexOf(chip);

    if (index >= 0) {
      this.values.splice(index, 1);
      this.propagateChange(this.values);
    }
  }

  reset() {
    this.chipInput.value = '';
    this.formControl.reset();
  }

  onSubmit(event) {
    if (this.chipInput.value === '' && !this.autoComplete.isOpen) {
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
