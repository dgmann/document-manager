import {ChangeDetectionStrategy, Component, OnDestroy, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {Record} from '@app/core/records';
import {InboxService} from './inbox.service';
import {ActionBarService} from '@app/inbox/action-bar';
import {untilDestroyed} from 'ngx-take-until-destroy';

@Component({
  selector: 'app-inbox',
  templateUrl: './inbox.component.html',
  styleUrls: ['./inbox.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InboxComponent implements OnInit, OnDestroy {
  records$: Observable<Record[]>;
  selectedRecord$: Observable<Record>;
  selectedIds: Observable<string[]>;
  isLoading$: Observable<boolean>;

  isDragInProgress = false;

  constructor(private inboxService: InboxService, private actionBar: ActionBarService) {
    this.isLoading$ = this.inboxService.isLoading$;
  }

  ngOnInit() {
    this.inboxService.loadRecords();
    this.records$ = this.inboxService.allInboxRecords$;
    this.selectedRecord$ = this.inboxService.selectedRecords$
      .pipe(map(records => records && records[0] || undefined));
    this.selectedIds = this.inboxService.selectedIds$;
    this.inboxService.selectedIds$
      .pipe(map(ids => ids.length > 1), untilDestroyed(this))
      .subscribe(isMultiselect => {
        if (isMultiselect) {
          if (!this.actionBar.isOpen) {
            this.actionBar.open();
          }
        } else {
          this.actionBar.dismiss();
        }
      });
  }

  ngOnDestroy(): void {
  }

  onSelectRecords(ids: string[]) {
    this.inboxService.selectIds(ids);
  }

  onDrop(event: DragEvent) {
    if (event.dataTransfer) {
      this.isDragInProgress = false;
      const files = event.dataTransfer.files;
      Array.from(files).forEach(file => this.inboxService.upload(file));
      this.preventAll(event);
      event.dataTransfer.clearData();
    }
  }

  onDragEnter(event: DragEvent) {
    const element = event.target as HTMLElement;
    if (element.id === 'inbox-document-list') {
      this.isDragInProgress = true;
    }
  }

  onDragLeave(event: DragEvent) {
    const element = event.target as HTMLElement;
    if (element.className === 'upload-symbol') {
      this.isDragInProgress = false;
    }
  }

  preventAll(event: DragEvent) {
    if (event.dataTransfer) {
      event.preventDefault();
      event.stopPropagation();
    }
  }
}
