import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs";
import {map} from "rxjs/operators";
import {PhysicianService} from "../physician.service";

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent implements OnInit {

  public reviewCount: Observable<number>;
  public escalatedCount: Observable<number>;
  public otherCount: Observable<number>;

  constructor(private physicianService: PhysicianService) {
    const countMap = map((records: any[]) => records.length);
    this.reviewCount = physicianService.getToReview().pipe(countMap);
    this.escalatedCount = physicianService.getEscalated().pipe(countMap);
    this.otherCount = physicianService.getOther().pipe(countMap);
  }

  ngOnInit() {
  }

}
