import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DocumentListComponent } from './document-list.component';
import { MatTableModule } from "@angular/material/table";
import { MatSortModule } from "@angular/material/sort";

@NgModule({
  imports: [
    CommonModule,
    MatSortModule,
    MatTableModule,
  ],
  declarations: [DocumentListComponent],
  exports: [
    DocumentListComponent
  ]
})
export class DocumentListModule { }
