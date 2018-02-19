import {Component} from '@angular/core';
import {MatDialog} from "@angular/material";
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs/Observable";
import {distinctUntilChanged, switchMap} from "rxjs/operators";
import {ReplaySubject} from "rxjs/ReplaySubject";
import {DocumentEditDialogComponent} from "../shared/document-edit-dialog/document-edit-dialog.component";
import {Record, RecordService, RequiredAction} from "../store";
import {InboxService} from "./inbox.service";

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss']
})
export class InboxComponent {

  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedRecordId = new ReplaySubject<string>();

  constructor(private inboxService: InboxService, private recordService: RecordService, public dialog: MatDialog) {
    inboxService.load();
    this.data = inboxService.all();
    const find = switchMap((id: string) => this.inboxService.find(id));
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

  upload(event: DropEvent) {
    for (let file of event.nativeEvent.dataTransfer.files) {
      this.recordService.upload(file)
    }
  }

  setRequiredAction(data: { record: Record, action: RequiredAction }) {
    this.recordService.update(data.record.id, {requiredAction: data.action})
  }
}
