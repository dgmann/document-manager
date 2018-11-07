import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EditorComponent } from './editor.component';
import { MatButtonModule, MatCardModule, MatIconModule } from "@angular/material";
import { EditorRouterModule } from "./editor.routes";
import { SharedModule } from "../shared";
import { EditorGuard } from "./editor.guard";
import { PageListComponent } from './page-list/page-list.component';
import { DragDropModule } from "@angular/cdk/drag-drop";
import { FlexLayoutModule } from "@angular/flex-layout";

@NgModule({
  imports: [
    CommonModule,
    FlexLayoutModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    SharedModule,
    EditorRouterModule,
    DragDropModule
  ],
  declarations: [EditorComponent, PageListComponent],
  exports: [EditorComponent],
  providers: [
    EditorGuard
  ]
})
export class EditorModule {
}
