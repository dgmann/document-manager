import { Component } from '@angular/core';
import { NotificationService } from "./shared/notification-service";
import { AutorefreshService } from "./store/record/autorefresh-service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  constructor(private autorefreshService: AutorefreshService,
              private notificationService: NotificationService,) {

    autorefreshService.start();
    this.notificationService.logToConsole();
    this.notificationService.logToSnackBar();
  }
}
