import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MatTableDataSource } from "@angular/material/table";
import sortBy from "lodash-es/sortBy";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { MultiSelectService } from "../multi-select.service";

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent implements OnInit {
  @Input() tags: Observable<string[]>;
  @Output() select = new EventEmitter<string[]>();
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();
  selectedTags: string[] = [];

  constructor(private multiselectService: MultiSelectService<string>) {
  }

  ngOnInit() {
    this.tags.pipe(
      map(tags => sortBy(tags))
    ).subscribe(data => this.dataSource.data = data);
  }

  onSelect(tag: string, event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.selectedTags = this.multiselectService.select(tag, this.dataSource.data, event);

    this.select.emit(this.selectedTags);
    return false;
  }

}
