import {Component, Inject} from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import {Record} from '@app/core/records';

@Component({
  selector: 'app-comment-dialog',
  templateUrl: './comment-dialog.component.html',
  styleUrls: ['./comment-dialog.component.scss']
})
export class CommentDialogComponent {
  public comment: string;

  constructor(public dialogRef: MatDialogRef<CommentDialogComponent>,
              @Inject(MAT_DIALOG_DATA) record: Record) {
    this.comment = record.comment;
  }

  onCloseClick() {
    this.dialogRef.close();
  }
}
