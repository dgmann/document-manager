import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from "@ngrx/store";
import { reducers } from "./reducers";

@NgModule({
  imports: [
    CommonModule,
    StoreModule.forFeature('records', reducers),
  ],
  declarations: []
})
export class CoreModule { }
