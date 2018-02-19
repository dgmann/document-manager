import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Record} from "../../store";
import {PhysicianService} from "../physician.service";

@Component({
  selector: 'app-page-review',
  templateUrl: './page-review.component.html',
  styleUrls: ['./page-review.component.scss']
})
export class PageReviewComponent implements OnInit {

  public records: Observable<Record[]>;

  constructor(private physicianService: PhysicianService) {
    this.records = physicianService.getToReview();
  }

  ngOnInit() {
  }

}
