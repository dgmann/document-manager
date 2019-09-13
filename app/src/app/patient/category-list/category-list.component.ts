import {Component, EventEmitter, Input, Output} from '@angular/core';
import {MatTableDataSource} from '@angular/material/table';
import {sortBy} from 'lodash-es';
import {Category} from '@app/core';
import {MultiSelectService} from '../multi-select.service';

@Component({
  selector: 'app-category-list',
  templateUrl: './category-list.component.html',
  styleUrls: ['./category-list.component.scss']
})
export class CategoryListComponent {
  @Output() selectCategories = new EventEmitter<Category[]>();

  @Input() set categories(categories: Category[]) {
    this.dataSource.data = sortBy(categories, ['name']);
  }
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<Category>();
  selectedCategories: Category[] = [];

  constructor(private multiselectService: MultiSelectService<Category>) {
  }

  onSelect(event: MouseEvent, category: Category) {
    event.preventDefault();
    event.stopPropagation();
    this.selectedCategories = this.multiselectService.select(category, this.dataSource.data, event);

    this.selectCategories.emit(this.selectedCategories);
    return false;
  }
}
