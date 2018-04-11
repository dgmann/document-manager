import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EditorComponent } from './editor.component';
import { MatCardModule } from "@angular/material";
import { EditorRouterModule } from "./editor.routes";
import { SharedModule } from "../shared";
import { DndModule } from "ng2-dnd";
import { EditorGuard } from "./editor.guard";

@NgModule({
  imports: [
    CommonModule,
    MatCardModule,
    SharedModule,
    DndModule,
    EditorRouterModule
  ],
  declarations: [EditorComponent],
  exports: [EditorComponent],
  providers: [
    EditorGuard
  ]
})
export class EditorModule {
}
