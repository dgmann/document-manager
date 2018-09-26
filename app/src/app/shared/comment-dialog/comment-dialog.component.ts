import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from "@angular/material";
import { Record } from "../../core/store/index";

@Component({
  selector: 'comment-dialog',
  templateUrl: './comment-dialog.component.html',
  styleUrls: ['./comment-dialog.component.css']
})
export class CommentDialogComponent {
  public comment: string;

  constructor(@Inject(MAT_DIALOG_DATA) record: Record) {
    this.comment = record.comment;
  }
}
