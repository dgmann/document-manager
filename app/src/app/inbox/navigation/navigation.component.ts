import {Component} from '@angular/core';
import {MatSlideToggleChange} from "@angular/material";
import {InboxService} from "../inbox.service";

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {

  constructor(private inboxService: InboxService) {
  }

  onChangeMultiSelect(event: MatSlideToggleChange) {
    this.inboxService.setMultiselect(event.checked);
  }

}
