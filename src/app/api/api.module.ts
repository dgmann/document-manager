import { NgModule } from '@angular/core';
import { HttpClientModule } from "@angular/common/http";
import { NgrxJsonApiModule } from "ngrx-json-api";


import { resourceDefinitions } from './resource-definitions';
import { RecordService } from "./record.service";
import { StoreModule } from "@ngrx/store";
import { EffectsModule } from "@ngrx/effects";

@NgModule({
  imports: [
    HttpClientModule,
    StoreModule.forFeature('api',{}),
    EffectsModule.forFeature([]),
    NgrxJsonApiModule.configure({
      apiUrl: 'http://localhost:8080',
      resourceDefinitions: resourceDefinitions,
    }),
  ],
  declarations: [],
  providers: [RecordService],
})
export class ApiModule { }
