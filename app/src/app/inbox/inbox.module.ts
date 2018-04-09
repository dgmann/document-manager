import {NgModule} from '@angular/core';
import {MatSlideToggleModule} from "@angular/material";
import {EffectsModule} from "@ngrx/effects";
import {StoreModule} from "@ngrx/store";
import {NgDragDropModule} from "ng-drag-drop";
import {SharedModule} from "../shared";
import {InboxComponent} from './inbox.component';
import {InboxRouterModule} from "./inbox.routes";
import {InboxService} from "./inbox.service";
import {NavigationComponent} from "./navigation/navigation.component";
import {metaReducers, reducers} from './reducers';
import {InboxEffects} from "./store/inbox.effects";

@NgModule({
  imports: [
    StoreModule.forFeature("inbox", reducers, {metaReducers}),
    EffectsModule.forFeature([InboxEffects]),
    InboxRouterModule,
    NgDragDropModule,
    MatSlideToggleModule,
    SharedModule
  ],
  declarations: [
    InboxComponent,
    NavigationComponent
  ],
  exports: [
    InboxComponent
  ],
  providers: [
    InboxService
  ]
})
export class InboxModule {
}
