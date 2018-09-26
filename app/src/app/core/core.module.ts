import { NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';
import { throwIfAlreadyLoaded } from "./module-import-guard";
import { NgDragDropModule } from "ng-drag-drop";
import { DndModule } from "ng2-dnd";
import { StoreModule } from "./store";

@NgModule({
  imports: [
    CommonModule,
    StoreModule,
    NgDragDropModule.forRoot(),
    DndModule.forRoot(),
  ],
  declarations: []
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    throwIfAlreadyLoaded(parentModule, 'CoreModule');
  }
}
