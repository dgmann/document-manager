import { Component } from '@angular/core';


import { Record } from "./core/record";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  constructor() {
  }

  DATA: Record[] = [
    {id: "1", date: new Date(2017, 12, 7), comment: "", type: "Scan"},
    {id: "2", date: new Date(2017, 12, 8), comment: "", type: "Fax"},
    {id: "3", date: new Date(2017, 12, 9), comment: "Neu?", type: "Fax"},
    {id: "4", date: new Date(2017, 12, 10), comment: "", type: "Scan"},
    {id: "5", date: new Date(2017, 12, 11), comment: "", type: "Scan"},
  ];
}
