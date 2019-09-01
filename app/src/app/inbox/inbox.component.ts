import {ChangeDetectionStrategy, Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {map, take} from 'rxjs/operators';
import {Record, Status} from '../core/store';
import {InboxService} from './inbox.service';

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InboxComponent implements OnInit {
  records: Observable<Record[]>;
  selectedRecord: Observable<Record>;
  selectedIds: Observable<string[]>;
  isMultiselect: Observable<boolean>;
  isLoading$: Observable<boolean>;

  constructor(private inboxService: InboxService) {
    this.isLoading$ = this.inboxService.isLoading$;
  }

  ngOnInit() {
    this.inboxService.loadRecords();
    this.records = this.inboxService.allInboxRecords$;
    this.selectedRecord = this.inboxService.selectedRecords$
      .pipe(map(records => records && records[0] || undefined));
    this.selectedIds = this.inboxService.selectedIds$;
    this.isMultiselect = this.inboxService.selectedIds$.pipe(map(ids => ids.length > 1));
  }

  onSelectRecords(ids: string[]) {
    this.inboxService.selectIds(ids);
  }

  onDrop(event: DragEvent) {
    if (event.dataTransfer) {
      const files = event.dataTransfer.files;
      for (let i = 0; i < files.length; i++) {
        this.inboxService.upload(files[i]);
      }
      this.preventAll(event);
      event.dataTransfer.clearData();
    }
  }

  preventAll(event: DragEvent) {
    if (event.dataTransfer) {
      event.preventDefault();
      event.stopPropagation();
    }
  }

  onSelectAllRecords(all: boolean) {
    if (all) {
      this.inboxService.allInboxRecordIds$
        .pipe(take(1))
        .subscribe(ids => this.inboxService.selectIds(ids));
    } else {
      this.inboxService.selectIds([]);
    }
  }

  onDeleteSelectedRecords() {
    this.inboxService.deleteSelectedRecords();
  }

  onSetStatusOfSelectedRecords(status: Status) {
    this.inboxService.updateSelectedRecords({status});
  }
}
