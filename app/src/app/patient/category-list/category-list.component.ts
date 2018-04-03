import {Component, OnInit} from '@angular/core';
import {MatTableDataSource} from "@angular/material";

@Component({
  selector: 'app-category-list',
  templateUrl: './category-list.component.html',
  styleUrls: ['./category-list.component.scss']
})
export class CategoryListComponent implements OnInit {
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();
  selectedTag = null;

  constructor() {
  }

  ngOnInit() {
  }

  select(tag: string) {
    if (this.selectedTag === tag) {
      this.selectedTag = null;
    } else {
      this.selectedTag = tag;
    }
  }

}
