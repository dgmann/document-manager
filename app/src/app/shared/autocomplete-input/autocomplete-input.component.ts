import {Component, Input, OnInit} from '@angular/core';
import {FormControl} from "@angular/forms";
import {Observable} from "rxjs/Observable";
import {concat, map, startWith, withLatestFrom} from "rxjs/operators";

@Component({
  selector: 'app-autocomplete-dropdown',
  templateUrl: './autocomplete-input.component.html',
  styleUrls: ['./autocomplete-input.component.scss']
})
export class AutocompleteInputComponent implements OnInit {
  @Input('options') options: Observable<string[]>;
  @Input('control') formControl: FormControl;

  filteredOptions: Observable<string[]>;

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
    return options.filter(option =>
      option.toLowerCase().indexOf(val.toLowerCase()) === 0);
  }
}
