import { Location } from "@angular/common";
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { includes } from 'lodash-es';
import { Observable } from "rxjs";
import { filter, switchMap, take } from "rxjs/operators";
import { Page, PageUpdate, Record, RecordService } from "../core/store";


@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss']
})
export class EditorComponent implements OnInit {
  record: Observable<Record>;
  pages: PageUpdate[];

  constructor(private recordService: RecordService,
              private route: ActivatedRoute,
              private location: Location) {
  }

  ngOnInit() {
    this.record = this.route.params.pipe(
      switchMap(params => this.recordService.find(params['id']))
    );
    this.record.subscribe(r => this.setPages(r.pages));
  }

  cancel() {
    this.location.back();
  }

  discard() {
    this.record.pipe(take(1)).subscribe(r => this.setPages(r.pages));
  }

  saveRecord() {
    this.record.pipe(take(1)).subscribe(r => {
      this.recordService.updatePages(r.id, this.pages);
      const sub = this.recordService.invalidIds.pipe(
        filter(ids => !includes(ids, r.id))
      ).subscribe(_ => {
        sub.unsubscribe();
        this.location.back();
      })
    });

  }

  setPages(pages: Page[]) {
    this.pages = pages.map(p => PageUpdate.FromPage(p))
  }

}
