import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {MatButtonModule, MatDialogModule, MatIconModule, MatTooltipModule} from "@angular/material";
import {MatSortModule} from "@angular/material/sort";
import {MatTableModule} from "@angular/material/table";
import {NgDragDropModule} from "ng-drag-drop";
import {DocumentListComponent} from './document-list.component';

@NgModule({
  imports: [
    CommonModule,
    MatSortModule,
    MatTableModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule,
    MatDialogModule,
    NgDragDropModule
  ],
  declarations: [DocumentListComponent],
  exports: [
    DocumentListComponent
  ]
})
export class DocumentListModule { }
