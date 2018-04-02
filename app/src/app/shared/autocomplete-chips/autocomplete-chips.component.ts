import {COMMA, ENTER} from "@angular/cdk/keycodes";
import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {FormControl} from "@angular/forms";
import {MatAutocompleteSelectedEvent, MatAutocompleteTrigger, MatChipInputEvent} from "@angular/material";
import {Observable} from "rxjs/Observable";
import {concat, map, startWith, withLatestFrom} from "rxjs/operators";

@Component({
  selector: 'app-autocomplete-chips',
  templateUrl: './autocomplete-chips.component.html',
  styleUrls: ['./autocomplete-chips.component.scss']
})
export class AutocompleteChipsComponent implements OnInit {
  @Input('value') value: string[];
  @Input('options') options: Observable<string[]>;
  @Input('control') formControl: FormControl;
  @ViewChild(MatAutocompleteTrigger) trigger: MatAutocompleteTrigger;

  filteredOptions: Observable<string[]>;
  separatorKeysCodes = [ENTER, COMMA];

  ngOnInit() {
    const filtered = this.formControl.valueChanges
      .pipe(
        startWith(''),
        withLatestFrom(this.options),
        map(([val, options]) => this.filter(val, options))
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

  add(event: MatChipInputEvent): void {
    let value = this.formControl.value;

    if ((value || '').trim()) {
      this.value.push(value.trim());
    }

    this.formControl.reset();
  }

  remove(chip: any): void {
    let index = this.value.indexOf(chip);

    if (index >= 0) {
      this.value.splice(index, 1);
    }
  }
}
