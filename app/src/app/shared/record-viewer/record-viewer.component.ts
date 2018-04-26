import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { Observable } from "rxjs/Observable";


import { Page, Record } from "../../store";

@Component({
  selector: 'app-record-viewer',
  templateUrl: './record-viewer.component.html',
  styleUrls: ['./record-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RecordViewerComponent {
  @Input('record') record: Observable<Record>;
  @Output('updatePages') updatePages = new EventEmitter<{ id: number, pages: Page[] }>();

  constructor() {
  }

  public up(recordId: number, pages: Page[], index: number) {
    if (index == 0) {
      return;
    }

    pages = pages.slice(0);
    const page = pages[index - 1];
    pages[index - 1] = pages[index];
    pages[index] = page;
    this.updatePages.emit({id: recordId, pages: pages});
  }
}
