import {CommonModule} from '@angular/common';
import {HttpClientModule} from "@angular/common/http";
import {NgModule} from '@angular/core';
import {EffectsModule} from '@ngrx/effects';
import {StoreModule as NgrxStore} from '@ngrx/store';
import {StoreDevtoolsModule} from '@ngrx/store-devtools';
import {environment} from '../../environments/environment';
import {AutorefreshService} from "./record/autorefresh-service";
import {RecordEffects} from './record/record.effects';
import {RecordService} from "./record/record.service";
import {metaReducers, reducers} from './reducers';

@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
    NgrxStore.forRoot(reducers, {metaReducers}),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
    EffectsModule.forRoot([RecordEffects])
  ],
  providers: [
    RecordService,
    AutorefreshService
  ],
  declarations: []
})
export class StoreModule {
}
