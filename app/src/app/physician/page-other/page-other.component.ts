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

  constructor(private physicianService: PhysicianService) {
    this.records = physicianService.getOther();
  }

  ngOnInit() {
  }

}
