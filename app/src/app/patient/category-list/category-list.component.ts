import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MatTableDataSource } from "@angular/material";
import { sortBy, uniq } from "lodash-es";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Category } from "../../core";

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

  onSelect(event: MouseEvent, category: Category) {
    event.preventDefault();
    event.stopPropagation();

    if (event.getModifierState("Control")) {
      let index = this.selectedCategories.indexOf(category);

      if (index >= 0) {
        this.selectedCategories.splice(index, 1);
      } else {
        this.selectedCategories = uniq([...this.selectedCategories, category]);
      }
    } else if (event.getModifierState("Shift")) {
      if (this.selectedCategories.length == 0) {
        this.selectedCategories = [category];
        return;
      }
      const firstCategory = this.selectedCategories[0];
      const selectFrom = this.dataSource.data.indexOf(firstCategory);
      const selectUntil = this.dataSource.data.indexOf(category);
      const indices = [selectFrom, selectUntil].sort();
      this.selectedCategories = this.dataSource.data.slice(indices[0], indices[1] + 1);
    } else {
      if (this.selectedCategories.length == 1 && this.selectedCategories[0] && this.selectedCategories[0] == category) {
        this.selectedCategories = []
      } else {
        this.selectedCategories = [category];
      }
    }

    this.select.emit(this.selectedCategories);
    return false;
  }
}
