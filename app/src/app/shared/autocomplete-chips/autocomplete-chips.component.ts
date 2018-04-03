import {COMMA, ENTER} from "@angular/cdk/keycodes";
import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {FormControl} from "@angular/forms";
import {MatAutocomplete, MatChipInputEvent, MatInput} from "@angular/material";
import {difference} from 'lodash-es';
import {Observable} from "rxjs/Observable";
import {concat, map, startWith, withLatestFrom} from "rxjs/operators";

@Component({
  selector: 'app-autocomplete-chips',
  templateUrl: './autocomplete-chips.component.html',
  styleUrls: ['./autocomplete-chips.component.scss']
})
export class AutocompleteChipsComponent implements OnInit {
  @Input('values') values: string[];
  @Input('options') options: Observable<string[]>;
  @Output('valuesChange') valuesChange = new EventEmitter<string[]>();
  @Output() submit = new EventEmitter();
  @ViewChild('chipInput') chipInput: MatInput;
  @ViewChild('auto') autoComplete: MatAutocomplete;
  formControl = new FormControl();

  filteredOptions: Observable<string[]>;
  separatorKeysCodes = [ENTER, COMMA];

  ngOnInit() {
    const filtered = this.formControl.valueChanges
      .pipe(
        startWith(''),
        withLatestFrom(this.options),
        map(([val, options]) => difference(this.filter(val, options), this.values))
      );
    this.filteredOptions = this.options.pipe(
      concat(filtered),
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
      this.valuesChange.emit(this.values);
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
      this.valuesChange.emit(this.values);
    }
  }

  reset() {
    this.chipInput['nativeElement'].value = '';
    this.formControl.reset();
  }

  onSubmit(event) {
    if (this.chipInput['nativeElement'].value == '' && !this.autoComplete.isOpen) {
      event.preventDefault();
      this.submit.emit(null);
    }
  }
}
