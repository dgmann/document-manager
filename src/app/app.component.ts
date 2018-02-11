import {Component} from '@angular/core';
import {MatDialog, MatSnackBar} from "@angular/material";
import {Observable} from "rxjs/Observable";
import {distinctUntilChanged, switchMap} from 'rxjs/operators';
import {ReplaySubject} from "rxjs/ReplaySubject";
import {DocumentEditDialogComponent} from "./document-edit-dialog/document-edit-dialog.component";
import {NotificationService} from "./shared/notification-service";


import {Record, RecordService} from "./store";
import {AutorefreshService} from "./store/record/autorefresh-service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedRecordId = new ReplaySubject<string>();

  constructor(private recordService: RecordService,
              private autorefreshService: AutorefreshService,
              private notificationService: NotificationService,
              public dialog: MatDialog,
              public snackbar: MatSnackBar) {
    recordService.load();
    this.data = recordService.all();
    const find = switchMap((id: string) => {
      let s = this.recordService.find(id);
      return s;
    });
    this.selectedRecord = this.selectedRecordId.pipe(distinctUntilChanged(), find);
    autorefreshService.start();
    this.notificationService.logToConsole();
    this.notificationService.logToSnackBar();
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
    });
  }

  deleteRecord(record: Record) {
    this.recordService.delete(record.id);
  }

  upload(event) {
    this.recordService.upload(event.mouseEvent.dataTransfer.files[0]);
  }
}
