import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';


import {PageUpdate, Record, RecordService} from '../../core/store/index';

@Component({
  selector: 'app-record-viewer',
  templateUrl: './record-viewer.component.html',
  styleUrls: ['./record-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RecordViewerComponent implements OnInit {
  @Input() showQuickEdit = true;
  viewMode = RecordViewerViewMode.MultipagePerRow;
  record: Record;
  pages: PageUpdate[];

  @Input('record') set setRecord(record: Record) {
    this.record = record;
    this.pages = record && record.pages.map(p => PageUpdate.FromPage(p)) || [];
  }

  constructor(private recordService: RecordService) {
  }

  public ngOnInit() {
  }

  public up(recordId: string, pages: PageUpdate[], index: number) {
    if (index === 0) {
      return;
    }

    pages = pages.slice(0);
    const page = pages[index - 1];
    pages[index - 1] = pages[index];
    pages[index] = page;
    this.recordService.updatePages(recordId, pages);
  }

  public down(recordId: string, pages: PageUpdate[], index: number) {
    if (index >= pages.length) {
      return;
    }

    pages = pages.slice(0);
    const page = pages[index + 1];
    pages[index + 1] = pages[index];
    pages[index] = page;
    this.recordService.updatePages(recordId, pages);
  }

  rotate(recordId: string, pages: PageUpdate[], index: number, degree: number) {
    pages = pages.slice(0);
    pages[index].rotate = this.mod(pages[index].rotate + degree, 360);
    this.recordService.updatePages(recordId, pages);
  }

  delete(recordId: string, pages: PageUpdate[], index: number) {
    pages = pages.slice(0);
    pages.splice(index, 1);
    this.recordService.updatePages(recordId, pages);
  }

  toogleViewMode() {
    if (this.viewMode === RecordViewerViewMode.MultipagePerRow) {
      this.viewMode = RecordViewerViewMode.SinglePagePerRow;
    } else {
      this.viewMode = RecordViewerViewMode.MultipagePerRow;
    }
  }

  mod(n, m) {
    return ((n % m) + m) % m;
  }

  trackByFn(index: number, item: PageUpdate) {
    return item.id;
  }
}

export enum RecordViewerViewMode {
  SinglePagePerRow = 'column',
  MultipagePerRow = 'row wrap'
}
