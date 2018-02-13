import { Component } from '@angular/core';
import { MatDialog, MatSnackBar } from "@angular/material";
import { Observable } from "rxjs/Observable";
import { distinctUntilChanged, switchMap } from 'rxjs/operators';
import { ReplaySubject } from "rxjs/ReplaySubject";
import { NotificationService } from "./shared/notification-service";


import { Record, RecordService } from "./store";
import { AutorefreshService } from "./store/record/autorefresh-service";

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
}
