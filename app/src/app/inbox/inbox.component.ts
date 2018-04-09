import {Component} from '@angular/core';
import {includes, without} from 'lodash-es';
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs/Observable";
import {map, take, withLatestFrom} from "rxjs/operators";
import {Record, RecordService} from "../store";
import {InboxService} from "./inbox.service";

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss']
})
export class InboxComponent {

  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedIds: Observable<string[]>;

  constructor(private inboxService: InboxService, private recordService: RecordService) {
    inboxService.load();
    this.data = inboxService.all();
    this.selectedRecord = this.inboxService.getSelectedRecords()
      .pipe(map(records => records && records[0] || undefined));
    this.selectedIds = this.inboxService.getSelectedIds();
  }

  selectRecord(record: Record) {
    this.inboxService.getSelectedIds()
      .pipe(
        take(1),
        withLatestFrom(this.inboxService.getMultiselect())
      )
      .subscribe(([ids, multiselect]) => {
        if (multiselect) {
          if (includes(ids, record.id)) {
            this.inboxService.selectIds(without(ids, record.id));
          } else {
            this.inboxService.selectIds([...ids, record.id]);
          }
        } else {
          this.inboxService.selectIds([record.id]);
        }
      });
  }

  updatePages(event) {
    this.recordService.update(event.id, {pages: event.pages})
  }

  upload(event: DropEvent) {
    for (let file of event.nativeEvent.dataTransfer.files) {
      this.recordService.upload(file)
    }
  }
}
