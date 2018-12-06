import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MatTableDataSource } from "@angular/material";
import { sortBy } from "lodash-es";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Category } from "../../core";
import { MultiSelectService } from "../multi-select.service";

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

  constructor(private multiselectService: MultiSelectService<Category>) {
  }

  ngOnInit() {
    this.categories.pipe(
      map(categories => sortBy(categories, ['name']))
    ).subscribe(data => this.dataSource.data = data);
  }

  onSelect(event: MouseEvent, category: Category) {
    event.preventDefault();
    event.stopPropagation();
    this.selectedCategories = this.multiselectService.select(category, this.dataSource.data, event);

    this.select.emit(this.selectedCategories);
    return false;
  }
}
