import {Component, OnInit} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Record} from "../../store";
import {PhysicianService} from "../physician.service";

@Component({
  selector: 'app-page-escalated',
  templateUrl: './page-escalated.component.html',
  styleUrls: ['./page-escalated.component.scss']
})
export class PageEscalatedComponent implements OnInit {

  public records: Observable<Record[]>;

  constructor(private physicianService: PhysicianService) {
    this.records = physicianService.getOther();
  }

  ngOnInit() {
  }

}
