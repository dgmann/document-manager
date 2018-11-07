import { Component, Input, OnInit } from '@angular/core';
import { PageUpdate } from "../../core/store";
import { CdkDragDrop, moveItemInArray } from "@angular/cdk/drag-drop";

@Component({
  selector: 'app-page-list',
  templateUrl: './page-list.component.html',
  styleUrls: ['./page-list.component.scss']
})
export class PageListComponent implements OnInit {
  @Input() pages: PageUpdate[];

  constructor() {
  }

  ngOnInit() {
  }

  drop(event: CdkDragDrop<string[]>) {
    moveItemInArray(this.pages, event.previousIndex, event.currentIndex);
  }

  rotate(page: PageUpdate, degree: number) {
    page.rotate = this.mod(page.rotate + degree, 360);
  }

  delete(page: PageUpdate) {
    const index = this.pages.indexOf(page);
    this.pages.splice(index, 1);
  }

  mod(n, m) {
    return ((n % m) + m) % m;
  }
}
