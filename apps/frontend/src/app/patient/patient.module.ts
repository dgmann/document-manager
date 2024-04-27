import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import { MatDateFnsModule } from '@angular/material-date-fns-adapter';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatTableModule} from '@angular/material/table';
import {MatDividerModule} from '@angular/material/divider';
import { RecordFilterModule } from '@app/patient/record-filter/record-filter.module';
import {EffectsModule} from '@ngrx/effects';
import {StoreModule} from '@ngrx/store';
import {SharedModule} from '../shared';
import {MultiRecordListComponent} from './multi-record-list/multi-record-list.component';
import {NavigationComponent} from './navigation/navigation.component';
import {PatientComponent} from './patient.component';
import {PatientRouterModule} from './patient.routes';
import {PatientService} from './patient.service';
import {metaReducers, reducers} from './reducers';
import {PatientEffects} from './store/patient.effects';
import {MatMenuModule} from '@angular/material/menu';

@NgModule({
  imports: [
    CommonModule,
    PatientRouterModule,
    SharedModule,
    MatCardModule,
    MatSidenavModule,
    MatIconModule,
    StoreModule.forFeature('patient', reducers, {metaReducers}),
    EffectsModule.forFeature([PatientEffects]),
    MatMenuModule,
    RecordFilterModule,
    MatButtonModule,
  ],
  declarations: [
    PatientComponent,
    NavigationComponent,
    MultiRecordListComponent
  ],
  providers: [
    PatientService
  ]
})
export class PatientModule {
}
