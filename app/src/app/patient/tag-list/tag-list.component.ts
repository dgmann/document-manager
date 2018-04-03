import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
import {uniq} from 'lodash-es';
import sortBy from "lodash-es/sortBy";
import {Observable} from "rxjs/Observable";
import {map, switchMap} from "rxjs/operators";
import {TagService} from "../../shared/tag-service";

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent implements OnInit {
  @Input('patientId') patientId: Observable<string>;
  @Output() selected = new EventEmitter<string[]>();
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();
  selectedTags: string[] = [];

  constructor(private tagService: TagService) {
  }

  ngOnInit() {
    this.patientId.pipe(
      switchMap(id => this.tagService.getByPatientId(id)),
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
