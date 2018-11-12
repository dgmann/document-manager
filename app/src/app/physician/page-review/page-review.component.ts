import { Component, OnInit } from '@angular/core';
import { Observable } from "rxjs";
import { Record } from "../../core/store/index";
import { PhysicianService } from "../physician.service";

@Component({
  selector: 'app-page-review',
  templateUrl: './page-review.component.html',
  styleUrls: ['./page-review.component.scss']
})
export class PageReviewComponent implements OnInit {
  public records: Observable<Record[]>;
  public selectedIds: Observable<string[]>;

  constructor(private physicianService: PhysicianService) {
  }

  ngOnInit() {
    this.records = this.physicianService.reviewRecords$;
    this.selectedIds = this.physicianService.selectedIds$;
  }

  selectRecord(record: Record) {
    this.physicianService.selectIds([record.id]);
  }
}
