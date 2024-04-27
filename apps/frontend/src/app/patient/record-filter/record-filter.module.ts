import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatDateFnsModule } from '@angular/material-date-fns-adapter';
import { MatCardModule } from '@angular/material/card';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatListModule } from '@angular/material/list';
import {
  DateRangeSelectorComponent
} from './date-range-selector/date-range-selector.component';
import { CategoryListComponent } from './category-list/category-list.component';
import { RecordFilterComponent } from './record-filter.component';
import { TagListComponent } from './tag-list/tag-list.component';



@NgModule({
  declarations: [
    RecordFilterComponent,
    CategoryListComponent,
    TagListComponent,
    DateRangeSelectorComponent,
  ],
  imports: [
    CommonModule,
    MatListModule,
    MatFormFieldModule,
    MatDatepickerModule,
    MatDateFnsModule,
    ReactiveFormsModule,
    MatCardModule
  ],
  exports: [
    RecordFilterComponent
  ]
})
export class RecordFilterModule { }
