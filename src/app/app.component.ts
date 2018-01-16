import { Component } from '@angular/core';
import { Observable } from "rxjs/Observable";


import { Record, RecordService } from "./api";
import { ReplaySubject } from "rxjs/ReplaySubject";
import { DocumentEditDialogComponent } from "./document-edit-dialog/document-edit-dialog.component";
import { MatDialog, MatSnackBar } from "@angular/material";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedRecordId = new ReplaySubject<string>();

  constructor(private recordService: RecordService, public dialog: MatDialog, public snackbar: MatSnackBar) {
    this.data = recordService.all();
    this.selectedRecord = this.selectedRecordId.distinctUntilChanged().switchMap(id => this.recordService.find(id));
  }

  selectRecord(record: Record) {
    this.selectedRecordId.next(record.id);
  }

  updatePages(event) {
    this.recordService.update(event.id, {pages: event.pages})
  }

  editRecord(record: Record) {
    this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record
    }).afterClosed().subscribe((result: Record) => {
      if (!result) {
        return;
      }
      this.recordService.update(result.id, {patientId: result.patientId, date: result.date, tags: result.tags});
      this.snackbar.open('Gespeichert', '', {
        duration: 3000
      });
    });
  }
}
