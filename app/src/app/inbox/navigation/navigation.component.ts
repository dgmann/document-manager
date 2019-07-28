import { Component, OnInit } from '@angular/core';
import { MatSlideToggleChange } from "@angular/material/slide-toggle";
import { Observable } from "rxjs";
import { InboxService } from "../inbox.service";

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent implements OnInit {
  isMultiSelect: Observable<boolean>;

  constructor(private inboxService: InboxService) {}

  ngOnInit() {
    this.isMultiSelect = this.inboxService.isMultiSelect$;
  }

  onChangeMultiSelect(event: MatSlideToggleChange) {
    this.inboxService.setMultiselect(event.checked);
  }

}
