import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Observable} from "rxjs";
import {Category, CategoryService} from "../../shared/category-service";
import {Record} from "../../store";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent implements OnInit {
  @Input('records') records: Observable<Record[]>;
  @Output('clickRecord') clickRecord = new EventEmitter<string>();
  categories: Observable<{ [id: string]: Category }>;

  constructor(private categoryService: CategoryService) {
  }

  ngOnInit() {
    this.categories = this.categoryService.getAsMap()
  }

  onRecordClicked(id: string) {
    this.clickRecord.emit(id);
  }

}
