import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-patient-search',
  templateUrl: './patient-search.component.html',
  styleUrls: ['./patient-search.component.scss']
})
export class PatientSearchComponent implements OnInit {

  public searchResults = [{
    name: 'John Doe',
    id: '1'
  }];

  constructor() {
  }

  ngOnInit() {
  }

  displayFn(patient): string | undefined {
    return patient ? patient.name : undefined;
  }

}
