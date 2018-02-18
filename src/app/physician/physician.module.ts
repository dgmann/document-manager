import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {NavigationComponent} from './navigation/navigation.component';
import {PhysicianComponent} from './physician.component';
import {PhysicianRouterModule} from "./physician.routes";

@NgModule({
  imports: [
    CommonModule,
    PhysicianRouterModule
  ],
  declarations: [PhysicianComponent, NavigationComponent]
})
export class PhysicianModule {
}
