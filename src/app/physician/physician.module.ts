import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {MatButtonModule} from "@angular/material";
import {NavigationComponent} from './navigation/navigation.component';
import {PhysicianComponent} from './physician.component';
import {PhysicianRouterModule} from "./physician.routes";

@NgModule({
  imports: [
    CommonModule,
    PhysicianRouterModule,
    MatButtonModule
  ],
  declarations: [PhysicianComponent, NavigationComponent]
})
export class PhysicianModule {
}
