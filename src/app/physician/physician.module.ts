import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {MatButtonModule} from "@angular/material";
import {StoreModule} from "@ngrx/store";
import {SharedModule} from "../shared";
import {NavigationComponent} from './navigation/navigation.component';
import {PageReviewComponent} from './page-review/page-review.component';
import {PhysicianComponent} from './physician.component';
import {PhysicianRouterModule} from "./physician.routes";
import {PhysicianService} from "./physician.service";
import {metaReducers, reducers} from "./reducers";

@NgModule({
  imports: [
    CommonModule,
    PhysicianRouterModule,
    StoreModule.forFeature("physician", reducers, {metaReducers}),
    MatButtonModule,
    SharedModule
  ],
  declarations: [
    PhysicianComponent,
    NavigationComponent,
    PageReviewComponent
  ],
  providers: [
    PhysicianService
  ]
})
export class PhysicianModule {
}
