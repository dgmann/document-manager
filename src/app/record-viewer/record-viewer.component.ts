import { Component, Input } from '@angular/core';
import { Observable } from "rxjs/Observable";


import { Page } from "../api";

@Component({
  selector: 'app-record-viewer',
  templateUrl: './record-viewer.component.html',
  styleUrls: ['./record-viewer.component.scss']
})
export class RecordViewerComponent {
  @Input('pages') pages: Observable<Page[]>;

  constructor() { }
}
