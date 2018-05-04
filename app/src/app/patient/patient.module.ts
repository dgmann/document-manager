import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import {
  MatAutocompleteModule,
  MatButtonModule,
  MatCardModule,
  MatDatepickerModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatSidenavModule,
  MatTableModule
} from "@angular/material";
import { MatDividerModule } from '@angular/material/divider';
import { EffectsModule } from "@ngrx/effects";
import { StoreModule } from "@ngrx/store";
import { SharedModule } from "../shared";
import { CategoryListComponent } from './category-list/category-list.component';
import { MultiRecordListComponent } from './multi-record-list/multi-record-list.component';
import { NavigationComponent } from './navigation/navigation.component';
import { PatientComponent } from './patient.component';
import { PatientRouterModule } from "./patient.routes";
import { PatientService } from "./patient.service";
import { metaReducers, reducers } from "./reducers";
import { PatientEffects } from "./store/patient.effects";
import { TagListComponent } from './tag-list/tag-list.component';
import { ThreeColumnPanelComponent } from './three-column-panel/three-column-panel.component';
import { RecordFilterComponent } from './record-filter/record-filter.component';
import { DateRangeSelectorComponent } from './date-range-selector/date-range-selector.component';
import { FormsModule } from "@angular/forms";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    PatientRouterModule,
    SharedModule,
    MatButtonModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    MatInputModule,
    MatTableModule,
    MatCardModule,
    MatDatepickerModule,
    MatDividerModule,
    MatSidenavModule,
    MatIconModule,
    StoreModule.forFeature("patient", reducers, {metaReducers}),
    EffectsModule.forFeature([PatientEffects]),
  ],
  declarations: [
    PatientComponent,
    NavigationComponent,
    TagListComponent,
    MultiRecordListComponent,
    CategoryListComponent,
    ThreeColumnPanelComponent,
    RecordFilterComponent,
    DateRangeSelectorComponent
  ],
  providers: [
    PatientService
  ]
})
export class PatientModule {
}
