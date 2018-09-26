import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { includes, without } from 'lodash-es';
import { DropEvent } from "ng-drag-drop";
import { Observable } from "rxjs";
import { map, take, withLatestFrom } from "rxjs/operators";
import { Record, Status } from "../core/store";
import { InboxService } from "./inbox.service";

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InboxComponent implements OnInit{
  records: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedIds: Observable<string[]>;
  isMultiselect: Observable<boolean>;

  constructor(private inboxService: InboxService) {
  }

  ngOnInit() {
    this.inboxService.loadRecords();
    this.records = this.inboxService.allInboxRecords$;
    this.selectedRecord = this.inboxService.selectedRecords$
      .pipe(map(records => records && records[0] || undefined));
    this.selectedIds = this.inboxService.selectedIds$;
    this.isMultiselect = this.inboxService.isMultiSelect$;
  }

  onSelectRecord(id: string) {
    this.inboxService.selectedIds$
      .pipe(
        take(1),
        withLatestFrom(this.inboxService.isMultiSelect$)
      )
      .subscribe(([ids, multiselect]) => {
        let idsToSelect = [];
        if (multiselect) {
          if (includes(ids, id)) {
            idsToSelect = without(ids, id);
          } else {
            idsToSelect = [...ids, id];
          }
        } else {
          idsToSelect = [id];
        }
        this.inboxService.selectIds(idsToSelect);
      });
  }

  onDrop(event: DropEvent) {
    for (let file of event.nativeEvent.dataTransfer.files) {
      this.inboxService.upload(file)
    }
  }

  onSelectAllRecords(all: boolean) {
    if (all) {
      this.inboxService.allInboxRecordIds$
        .pipe(take(1))
        .subscribe(ids => this.inboxService.selectIds(ids));
    }
    else {
      this.inboxService.selectIds([]);
    }
  }

  onDeleteSelectedRecords() {
    this.inboxService.deleteSelectedRecords();
  }

  onSetStatusOfSelectedRecords(status: Status) {
    this.inboxService.updateSelectedRecords({status})
  }
}
