import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { Observable } from "rxjs";


import { PageUpdate, Record, RecordService } from "../../store";
import { map } from "rxjs/operators";

@Component({
  selector: 'app-record-viewer',
  templateUrl: './record-viewer.component.html',
  styleUrls: ['./record-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RecordViewerComponent implements OnInit {
  @Input('record') record: Observable<Record>;

  pages$: Observable<PageUpdate[]>;

  constructor(private recordService: RecordService) {
  }

  public ngOnInit() {
    this.pages$ = this.record.pipe(map(r => r.pages.map(p => PageUpdate.FromPage(p))))
  }

  public up(recordId: string, pages: PageUpdate[], index: number) {
    if (index == 0) {
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

  mod(n, m) {
    return ((n % m) + m) % m;
  }

  trackByFn(index: number, item: PageUpdate) {
    return item.id;
  }
}
