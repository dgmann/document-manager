import { Component } from '@angular/core';
import { MatDialog } from "@angular/material";
import { Observable } from "rxjs/Observable";
import { distinctUntilChanged, switchMap } from "rxjs/operators";
import { ReplaySubject } from "rxjs/ReplaySubject";
import { DocumentEditDialogComponent } from "../document-edit-dialog/document-edit-dialog.component";
import { Record, RecordService } from "../store";

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss']
})
export class InboxComponent {

  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedRecordId = new ReplaySubject<string>();

  constructor(private recordService: RecordService, public dialog: MatDialog) {
    this.data = recordService.all();
    const find = switchMap((id: string) => {
      let s = this.recordService.find(id);
      return s;
    });
    this.selectedRecord = this.selectedRecordId.pipe(distinctUntilChanged(), find);
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

  appendRecord(event) {
    this.recordService.append(event.source.id, event.target.id);
  }
}
