import {ChangeDetectionStrategy, Component} from '@angular/core';
import {Title} from '@angular/platform-browser';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';
import {Observable} from 'rxjs';
import {filter, map, mergeMap} from 'rxjs/operators';
import {Patient} from './patient';
import {AutorefreshService} from '@app/core/autorefresh';
import {NotificationService} from '@app/core/notifications';
import { ExternalApiService } from './shared/document-edit-dialog/external-api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AppComponent {
  public title: Observable<string>;

  constructor(private autorefreshService: AutorefreshService,
              private notificationService: NotificationService,
              private router: Router,
              private activatedRoute: ActivatedRoute,
              private titleService: Title,
              private patientService: ExternalApiService) {

    autorefreshService.start();
    this.notificationService.logToConsole();
    this.notificationService.logToSnackBar();

    this.title = this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map(() => this.activatedRoute),
      map(route => {
        while (route.firstChild) {
          route = route.firstChild;
        }
        return route;
      }),
      filter(route => route.outlet === 'primary'),
      mergeMap(route => route.data),
      map(data => data.title)
    );
    this.title.subscribe(title => this.titleService.setTitle(`${title} - Document Manager`));
  }

  onSelectPatient(event: Patient) {
    this.router.navigate(['/patient', event.id]);
  }

  navigateToCurrentPatient() {
    this.patientService.getSelectedPatient().subscribe(patient => this.onSelectPatient(patient))
  }
}
