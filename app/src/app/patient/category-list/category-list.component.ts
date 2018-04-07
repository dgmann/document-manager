import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
import {uniq} from "lodash-es";
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
  @Output() select = new EventEmitter<Category[]>();
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<Category>();
  selectedCategories: Category[] = [];

  constructor() {
  }

  ngOnInit() {
    this.categories.pipe(
      map(categories => sortBy(categories, ['name']))
    ).subscribe(data => this.dataSource.data = data);
  }

  onSelect(category: Category) {
    let index = this.selectedCategories.indexOf(category);

    if (index >= 0) {
      this.selectedCategories.splice(index, 1);
    } else {
      this.selectedCategories = uniq([...this.selectedCategories, category]);
    }
    this.select.emit(this.selectedCategories);
  }

}
