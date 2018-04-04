import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
import {uniq} from 'lodash-es';
import sortBy from "lodash-es/sortBy";
import {Observable} from "rxjs/Observable";
import {map} from "rxjs/operators";

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent implements OnInit {
  @Input() tags: Observable<string[]>;
  @Output() selected = new EventEmitter<string[]>();
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();
  selectedTags: string[] = [];

  constructor() {
  }

  ngOnInit() {
    this.tags.pipe(
      map(tags => sortBy(tags))
    ).subscribe(data => this.dataSource.data = data);
  }

  select(tag: string) {
    let index = this.selectedTags.indexOf(tag);

    if (index >= 0) {
      this.selectedTags.splice(index, 1);
    } else {
      this.selectedTags = uniq([...this.selectedTags, tag]);
    }
    this.selected.emit(this.selectedTags);
  }

}
