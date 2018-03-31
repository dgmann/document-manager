import {Component, Input, OnInit} from '@angular/core';
import {MatTableDataSource} from "@angular/material";
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
  displayedColumns = ['main'];
  dataSource = new MatTableDataSource<string>();

  constructor(private tagService: TagService) {
  }

  ngOnInit() {
    this.patientId.pipe(
      switchMap(id => this.tagService.getByPatientId(id)),
      map(tags => sortBy(tags))
    ).subscribe(data => this.dataSource.data = data);
  }

}
