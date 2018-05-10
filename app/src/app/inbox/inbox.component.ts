import {ChangeDetectionStrategy, Component, OnInit} from '@angular/core';
import {includes, without} from 'lodash-es';
import {DropEvent} from "ng-drag-drop";
import {Observable} from "rxjs";
import {map, take, withLatestFrom} from "rxjs/operators";
import {Record, RecordService, Status} from "../store";
import {InboxService} from "./inbox.service";

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InboxComponent implements OnInit{
  data: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedIds: Observable<string[]>;
  isMultiselect: Observable<boolean>;

  constructor(private inboxService: InboxService,
              private recordService: RecordService) {
  }

  ngOnInit() {
    this.recordService.load({status: Status.INBOX});
    this.data = this.inboxService.all();
    this.selectedRecord = this.inboxService.getSelectedRecords()
      .pipe(map(records => records && records[0] || undefined));
    this.selectedIds = this.inboxService.getSelectedIds();
    this.isMultiselect = this.inboxService.getMultiselect();
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

  upload(event: DropEvent) {
    for (let file of event.nativeEvent.dataTransfer.files) {
      this.recordService.upload(file)
    }
  }

  selectAllRecords(all: boolean) {
    if (all) {
      this.data.pipe(take(1)).subscribe((records: Record[]) => this.inboxService.selectIds(records.map(r => r.id)));
    }
    else {
      this.inboxService.selectIds([]);
    }
  }
}
