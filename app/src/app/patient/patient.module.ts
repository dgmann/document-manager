import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {
  MatAutocompleteModule,
  MatButtonModule,
  MatCardModule,
  MatFormFieldModule,
  MatInputModule,
  MatTableModule
} from "@angular/material";
import {MatDividerModule} from '@angular/material/divider';
import {SharedModule} from "../shared";
import {MultiRecordListComponent} from './multi-record-list/multi-record-list.component';
import {NavigationComponent} from './navigation/navigation.component';
import {PatientComponent} from './patient.component';
import {PatientRouterModule} from "./patient.routes";
import {TagListComponent} from './tag-list/tag-list.component';

@NgModule({
  imports: [
    CommonModule,
    PatientRouterModule,
    SharedModule,
    MatButtonModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    MatInputModule,
    MatTableModule,
    MatCardModule,
    MatDividerModule
  ],
  declarations: [PatientComponent, NavigationComponent, TagListComponent, MultiRecordListComponent]
})
export class PatientModule {
}
