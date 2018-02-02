import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule as NgrxStore } from '@ngrx/store';
import { metaReducers, reducers } from './reducers';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from '../../environments/environment';
import { EffectsModule } from '@ngrx/effects';
import { RecordEffects } from './record/record.effects';
import { HttpClientModule } from "@angular/common/http";
import { RecordService } from "./record/record.service";

@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
    NgrxStore.forRoot(reducers, {metaReducers}),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
    EffectsModule.forRoot([RecordEffects])
  ],
  providers: [
    RecordService
  ],
  declarations: []
})
export class StoreModule {
}
