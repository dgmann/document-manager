import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {MatAutocompleteModule, MatButtonModule, MatFormFieldModule, MatInputModule} from "@angular/material";
import {SharedModule} from "../shared";
import {NavigationComponent} from './navigation/navigation.component';
import {PatientComponent} from './patient.component';
import {PatientRouterModule} from "./patient.routes";

@NgModule({
  imports: [
    CommonModule,
    PatientRouterModule,
    SharedModule,
    MatButtonModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    MatInputModule
  ],
  declarations: [PatientComponent, NavigationComponent]
})
export class PatientModule {
}
