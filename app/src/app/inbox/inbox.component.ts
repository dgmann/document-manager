import {Component} from '@angular/core';
import {MatDialog} from "@angular/material";
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs/Observable";
import {map} from "rxjs/operators";
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

  constructor(private inboxService: InboxService, private recordService: RecordService, public dialog: MatDialog) {
    inboxService.load();
    this.data = inboxService.all();
    this.selectedRecord = this.inboxService.getSelectedRecords()
      .pipe(map(records => records && records[0] || undefined))
  }

  selectRecord(record: Record) {
    this.inboxService.selectIds([record.id]);
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
      this.recordService.update(result.id, {
        patientId: result.patientId,
        date: result.date,
        tags: result.tags,
        categoryId: result.categoryId
      });
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
