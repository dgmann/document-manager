import {ChangeDetectionStrategy, Component, Inject, OnInit} from '@angular/core';
import {MAT_SNACK_BAR_DATA} from '@angular/material/snack-bar';
import {NotificationEvent} from '@app/core/notifications';

@Component({
  selector: 'app-event-snackbar',
  templateUrl: './event-snackbar.component.html',
  styleUrls: ['./event-snackbar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class EventSnackbarComponent implements OnInit {
  constructor(@Inject(MAT_SNACK_BAR_DATA) public events: NotificationEvent[]) {
  }

  ngOnInit() {
  }
}
