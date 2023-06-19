import {Component, EventEmitter, Input, Output} from '@angular/core';
import { MatSelectionListChange } from '@angular/material/list';
import sortBy from 'lodash-es/sortBy';

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent {
  @Output() selectTag = new EventEmitter<string[]>();

  @Input() set tags(tags: string[]) {
    this.dataSource = sortBy(tags);
  }
  dataSource: string[] = [];

  constructor() {
  }

  onSelectionChange($event: MatSelectionListChange) {
    const selected = $event.options.map(option => option.value);
    this.selectTag.emit(selected);
  }
}
