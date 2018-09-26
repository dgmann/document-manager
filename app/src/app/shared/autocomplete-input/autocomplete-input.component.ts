import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { FormControl } from "@angular/forms";
import { concat, Observable } from "rxjs";
import { map, startWith, withLatestFrom } from "rxjs/operators";
import { Category } from "../../core";

@Component({
  selector: 'app-autocomplete-dropdown',
  templateUrl: './autocomplete-input.component.html',
  styleUrls: ['./autocomplete-input.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
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
    this.filteredOptions = concat(this.options, filtered);
  }

  filter(val: string, options: Category[]): Category[] {
    return options.filter(option =>
      option.name.toLowerCase().indexOf(val.toLowerCase()) === 0);
  }

  compareByValue(f1: any, f2: any) {
    return f1 && f2 && f1.id === f2.id;
  }
}
