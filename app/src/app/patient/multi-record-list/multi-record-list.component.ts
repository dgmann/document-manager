import {Component, Input, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Category, CategoryService} from "../../shared/category-service";
import {Record} from "../../store";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss']
})
export class MultiRecordListComponent implements OnInit {
  @Input('records') records: Observable<Record[]>;
  categories: { [id: string]: Category } = {};

  constructor(private categoryService: CategoryService) {
  }

  ngOnInit() {
    this.categoryService.get().subscribe(cat => cat.forEach(c => this.categories[c.id] = c));
  }

}
