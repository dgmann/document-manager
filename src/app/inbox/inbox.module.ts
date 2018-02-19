import {NgModule} from '@angular/core';
import {EffectsModule} from "@ngrx/effects";
import {StoreModule} from "@ngrx/store";
import {SharedModule} from "../shared";
import {InboxComponent} from './inbox.component';
import {InboxRouterModule} from "./inbox.routes";
import {InboxService} from "./inbox.service";
import {metaReducers, reducers} from './reducers';
import {InboxEffects} from "./store/inbox.effects";

@NgModule({
  imports: [
    StoreModule.forFeature("inbox", reducers, {metaReducers}),
    EffectsModule.forFeature([InboxEffects]),
    InboxRouterModule,
    SharedModule
  ],
  declarations: [
    InboxComponent,
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
