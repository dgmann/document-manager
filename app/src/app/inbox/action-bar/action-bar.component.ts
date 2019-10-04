import {ChangeDetectionStrategy, Component, OnInit} from '@angular/core';
import {Status} from '@app/core/store';
import {InboxService} from '@app/inbox/inbox.service';
import {take} from 'rxjs/operators';

@Component({
  selector: 'app-action-bar',
  templateUrl: './action-bar.component.html',
  styleUrls: ['./action-bar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ActionBarComponent implements OnInit {
  status = Status;

  constructor(private inboxService: InboxService) {
  }

  ngOnInit() {
  }

  onDeleteRecord(event) {
    event.stopPropagation();
    this.inboxService.deleteSelectedRecords();
  }

  setStatus(status: Status) {
    this.inboxService.updateSelectedRecords({status});
  }

  selectAll(all: boolean) {
    if (all) {
      this.inboxService.allInboxRecordIds$
        .pipe(take(1))
        .subscribe(ids => this.inboxService.selectIds(ids));
    } else {
      this.inboxService.selectIds([]);
    }
  }
}
