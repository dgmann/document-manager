import { Component, Input, OnInit } from '@angular/core';
import { Page } from "../../store";

@Component({
  selector: 'app-page-list',
  templateUrl: './page-list.component.html',
  styleUrls: ['./page-list.component.scss']
})
export class PageListComponent implements OnInit {
  @Input() pages: Page[];

  constructor() {
  }

  ngOnInit() {
  }

}
