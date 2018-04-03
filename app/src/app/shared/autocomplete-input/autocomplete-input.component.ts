import {Component, Input, OnInit} from '@angular/core';
import {FormControl} from "@angular/forms";
import {Observable} from "rxjs/Observable";
import {concat, map, startWith, withLatestFrom} from "rxjs/operators";
import {Category} from "../category-service";

@Component({
  selector: 'app-autocomplete-dropdown',
  templateUrl: './autocomplete-input.component.html',
  styleUrls: ['./autocomplete-input.component.scss']
})
export class AutocompleteInputComponent implements OnInit {
  @Input('options') options: Observable<Category[]>;
  @Input('control') formControl: FormControl;

  filterFormControl = new FormControl();
  filteredOptions: Observable<Category[]>;

  ngOnInit() {
    const filtered = this.filterFormControl.valueChanges
      .pipe(
        startWith(''),
        withLatestFrom(this.options),
        map(([val, options]) => this.filter(val, options))
      );
    this.filteredOptions = this.options.pipe(
      concat(filtered),
    );
  }

  filter(val: string, options: Category[]): Category[] {
    return options.filter(option =>
      option.name.toLowerCase().indexOf(val.toLowerCase()) === 0);
  }

  displayFn(category?: Category): string | undefined {
    return category ? category.name : undefined;
  }
}
