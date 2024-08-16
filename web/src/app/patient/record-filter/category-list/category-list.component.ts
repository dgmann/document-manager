import {Component, EventEmitter, Input, Output} from '@angular/core';
import { MatSelectionListChange } from '@angular/material/list';
import {sortBy} from 'lodash-es';
import {Category} from '@app/core/categories';

@Component({
  selector: 'app-category-list',
  templateUrl: './category-list.component.html',
  styleUrls: ['./category-list.component.scss']
})
export class CategoryListComponent {
  @Output() selectCategories = new EventEmitter<Category[]>();

  @Input() set categories(categories: Category[]) {
    this.dataSource = sortBy(categories, ['name']);
  }
  dataSource: Category[] = [];

  constructor() {
  }

  onSelectionChange($event: MatSelectionListChange) {
    const selected = $event.source.selectedOptions.selected.map(option => option.value);
    this.selectCategories.emit(selected);
  }
}
