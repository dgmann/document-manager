import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from "@angular/material";
import { EffectsModule } from "@ngrx/effects";
import { StoreModule } from "@ngrx/store";
import { SharedModule } from "../shared";
import { NavigationComponent } from './navigation/navigation.component';
import { PageEscalatedComponent } from './page-escalated/page-escalated.component';
import { PageOtherComponent } from './page-other/page-other.component';
import { PageReviewComponent } from './page-review/page-review.component';
import { PhysicianComponent } from './physician.component';
import { PhysicianRouterModule } from "./physician.routes";
import { metaReducers, reducers } from "./reducers";
import { PhysicianEffects } from "./store/physician.effects";

@NgModule({
  imports: [
    CommonModule,
    PhysicianRouterModule,
    StoreModule.forFeature("physician", reducers, {metaReducers}),
    EffectsModule.forFeature([PhysicianEffects]),
    MatButtonModule,
    SharedModule
  ],
  declarations: [
    PhysicianComponent,
    NavigationComponent,
    PageReviewComponent,
    PageEscalatedComponent,
    PageOtherComponent
  ]
})
export class PhysicianModule {
}
