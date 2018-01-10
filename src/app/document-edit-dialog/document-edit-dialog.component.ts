import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatChipInputEvent } from "@angular/material";
import { COMMA, ENTER } from "@angular/cdk/keycodes";


import { Record } from "../api";


@Component({
  selector: 'app-document-edit-dialog',
  templateUrl: './document-edit-dialog.component.html',
  styleUrls: ['./document-edit-dialog.component.scss']
})
export class DocumentEditDialogComponent {
  selectable: boolean = true;
  removable: boolean = true;
  addOnBlur: boolean = true;

  // Enter, comma
  separatorKeysCodes = [ENTER, COMMA];

  constructor(@Inject(MAT_DIALOG_DATA) public record: Record) {
  }

  add(event: MatChipInputEvent): void {
  }

  remove(fruit: any): void {
  }
}
