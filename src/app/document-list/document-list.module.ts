import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DocumentListComponent } from './document-list.component';
import { MatTableModule } from "@angular/material/table";
import { MatSortModule } from "@angular/material/sort";
import { MatButtonModule, MatIconModule, MatTooltipModule } from "@angular/material";

@NgModule({
  imports: [
    CommonModule,
    MatSortModule,
    MatTableModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule
  ],
  declarations: [DocumentListComponent],
  exports: [
    DocumentListComponent
  ]
})
export class DocumentListModule { }
