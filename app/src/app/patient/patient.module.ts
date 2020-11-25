import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatTableModule} from '@angular/material/table';
import {MatTabsModule} from '@angular/material/tabs';
import {MatDividerModule} from '@angular/material/divider';
import {EffectsModule} from '@ngrx/effects';
import {StoreModule} from '@ngrx/store';
import {SharedModule} from '../shared';
import {CategoryListComponent} from './category-list/category-list.component';
import {MultiRecordListComponent} from './multi-record-list/multi-record-list.component';
import {NavigationComponent} from './navigation/navigation.component';
import {PatientComponent} from './patient.component';
import {PatientRouterModule} from './patient.routes';
import {PatientService} from './patient.service';
import {metaReducers, reducers} from './reducers';
import {PatientEffects} from './store/patient.effects';
import {TagListComponent} from './tag-list/tag-list.component';
import {RecordFilterComponent} from './record-filter/record-filter.component';
import {DateRangeSelectorComponent} from './date-range-selector/date-range-selector.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatMenuModule} from '@angular/material/menu';
import { MatMomentDateModule } from '@angular/material-moment-adapter';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    PatientRouterModule,
    SharedModule,
    MatButtonModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    MatInputModule,
    MatTableModule,
    MatCardModule,
    MatDatepickerModule,
    MatMomentDateModule,
    MatDividerModule,
    MatSidenavModule,
    MatIconModule,
    MatTabsModule,
    StoreModule.forFeature('patient', reducers, {metaReducers}),
    EffectsModule.forFeature([PatientEffects]),
    MatMenuModule,
  ],
  declarations: [
    PatientComponent,
    NavigationComponent,
    TagListComponent,
    MultiRecordListComponent,
    CategoryListComponent,
    RecordFilterComponent,
    DateRangeSelectorComponent
  ],
  providers: [
    PatientService
  ]
})
export class PatientModule {
}
