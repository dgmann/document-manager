import { Component, OnInit } from '@angular/core';
import { Page, Record, RecordService } from "../store";
import { Observable } from "rxjs/Observable";
import { ActivatedRoute } from "@angular/router";
import { switchMap } from "rxjs/operators";

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss']
})
export class EditorComponent implements OnInit {
  record: Observable<Record>;
  pages: Page[];

  constructor(private recordService: RecordService,
              private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.record = this.route.params.pipe(switchMap(params => this.recordService.find(params['id'])));
    this.record.subscribe(r => this.pages = r.pages);
  }

}
