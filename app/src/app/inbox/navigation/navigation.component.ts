import {Component, OnInit} from '@angular/core';
import {MatSlideToggleChange} from "@angular/material";
import {Observable} from "rxjs";
import {InboxService} from "../inbox.service";

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent implements OnInit {
  isMultiSelect: Observable<boolean>;

  constructor(private inboxService: InboxService) {}

  ngOnInit() {
    this.isMultiSelect = this.inboxService.getMultiselect();
  }

  onChangeMultiSelect(event: MatSlideToggleChange) {
    this.inboxService.setMultiselect(event.checked);
  }

}
