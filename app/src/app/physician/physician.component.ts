import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { Observable } from "rxjs";
import { filter, map } from "rxjs/operators";
import { Record, RecordService, RequiredAction } from "../store";
import { PhysicianService } from "./physician.service";

@Component({
  selector: 'app-physician',
  templateUrl: './physician.component.html',
  styleUrls: ['./physician.component.scss']
})
export class PhysicianComponent implements OnInit {
  selectedRecord: Observable<Record>;

  constructor(private physicianService: PhysicianService,
              private recordService: RecordService,
              private route: ActivatedRoute,
              private router: Router) {
  }

  ngOnInit() {
    this.selectedRecord = this.physicianService.getSelectedRecords().pipe(
      filter(records => records.length > 0),
      map(records => records[0])
    );
    this.physicianService.load(RequiredAction.REVIEW);
    this.physicianService.load(RequiredAction.ESCALATED);
    this.physicianService.load(RequiredAction.OTHER);

    this.route
      .queryParams
      .subscribe(params => {
        const id = params['selected'] || null;
        this.physicianService.selectIds([id]);
      });

    this.physicianService.getSelectedIds().subscribe(ids => this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        ...this.route.snapshot.queryParams,
        selected: ids.length > 0 ? ids[0] : null,
      }
    }))
  }
}
