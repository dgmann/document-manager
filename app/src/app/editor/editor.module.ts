import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EditorComponent } from './editor.component';
import { MatButtonModule, MatCardModule, MatIconModule } from "@angular/material";
import { EditorRouterModule } from "./editor.routes";
import { SharedModule } from "../shared";
import { DndModule } from "ng2-dnd";
import { EditorGuard } from "./editor.guard";
import { PageListComponent } from './page-list/page-list.component';

@NgModule({
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    SharedModule,
    DndModule,
    EditorRouterModule
  ],
  declarations: [EditorComponent, PageListComponent],
  exports: [EditorComponent],
  providers: [
    EditorGuard
  ]
})
export class EditorModule {
}
