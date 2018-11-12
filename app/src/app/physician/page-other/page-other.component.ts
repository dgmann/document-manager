import { Component, OnInit } from '@angular/core';
import { Observable } from "rxjs";
import { Record } from "../../core/store/index";
import { PhysicianService } from "../physician.service";

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
    this.records = this.physicianService.otherRecords$;
    this.selectedIds = this.physicianService.selectedIds$;
  }

  selectRecord(record: Record) {
    this.physicianService.selectIds([record.id]);
  }
}
