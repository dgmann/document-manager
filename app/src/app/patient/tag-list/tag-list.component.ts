import {Component, EventEmitter, Input, Output} from '@angular/core';
import {MatTableDataSource} from '@angular/material/table';
import sortBy from 'lodash-es/sortBy';
import {MultiSelectService} from '../multi-select.service';

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent {
  @Output() selectTag = new EventEmitter<string[]>();

  @Input() set tags(tags: string[]) {
    this.dataSource.data = sortBy(tags);
  }

  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();
  selectedTags: string[] = [];

  constructor(private multiselectService: MultiSelectService<string>) {
  }

  onSelect(tag: string, event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.selectedTags = this.multiselectService.select(tag, this.dataSource.data, event);

    this.selectTag.emit(this.selectedTags);
    return false;
  }

}
