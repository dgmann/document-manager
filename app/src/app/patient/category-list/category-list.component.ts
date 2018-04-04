import {Component, Input, OnInit} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
import sortBy from "lodash-es/sortBy";
import {Observable} from "rxjs/Observable";
import {map} from "rxjs/operators";
import {Category} from "../../shared/category-service";

@Component({
  selector: 'app-category-list',
  templateUrl: './category-list.component.html',
  styleUrls: ['./category-list.component.scss']
})
export class CategoryListComponent implements OnInit {
  @Input() categories: Observable<Category[]>;
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<Category>();
  selectedCategory = null;

  constructor() {
  }

  ngOnInit() {
    this.categories.pipe(
      map(categories => sortBy(categories, ['name']))
    ).subscribe(data => this.dataSource.data = data);
  }

  select(category: string) {
    if (this.selectedCategory === category) {
      this.selectedCategory = null;
    } else {
      this.selectedCategory = category;
    }
  }

}
