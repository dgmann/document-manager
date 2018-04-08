import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Record} from "../../store";
import {PhysicianService} from "../physician.service";

@Component({
  selector: 'app-page-other',
  templateUrl: './page-other.component.html',
  styleUrls: ['./page-other.component.scss']
})
export class PageOtherComponent implements OnInit {
  public records: Observable<Record[]>;
  public selectedIds: Observable<string[]>;

  constructor(private physicianService: PhysicianService) {
  }

  ngOnInit() {
    this.records = this.physicianService.getOther();
    this.selectedIds = this.physicianService.getSelectedIds();
  }

  selectRecord(record: Record) {
    this.physicianService.selectIds([record.id]);
  }
}
