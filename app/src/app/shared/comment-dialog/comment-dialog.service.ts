import {Injectable} from '@angular/core';
import {MatDialog} from '@angular/material/dialog';
import {Record} from '@app/core/store';
import {Observable} from 'rxjs';
import {CommentDialogComponent} from './comment-dialog.component';

@Injectable({
  providedIn: 'root'
})
export class CommentDialogService {
  constructor(private dialog: MatDialog) {
  }

  open(record: Record): Observable<string> {
    return this.dialog.open(CommentDialogComponent, {
      disableClose: true,
      data: record,
      width: '635px'
    }).afterClosed();
  }
}
