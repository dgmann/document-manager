import { Component, OnInit } from '@angular/core';
import { Observable } from "rxjs";
import { Record } from "../../core/store/index";
import { PhysicianService } from "../physician.service";

@Component({
  selector: 'app-page-escalated',
  templateUrl: './page-escalated.component.html',
  styleUrls: ['./page-escalated.component.scss']
})
export class PageEscalatedComponent implements OnInit {
  public records: Observable<Record[]>;
  public selectedIds: Observable<string[]>;

  constructor(private physicianService: PhysicianService) {
  }

  ngOnInit() {
    this.records = this.physicianService.getEscalated();
    this.selectedIds = this.physicianService.getSelectedIds();
  }

  selectRecord(record: Record) {
    this.physicianService.selectIds([record.id]);
  }
}
