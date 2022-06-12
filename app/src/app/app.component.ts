import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import {Title} from '@angular/platform-browser';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';
import {PatientSearchComponent} from '@app/shared/patient-search/patient-search.component';
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
  public navColor: Observable<string>;

  constructor(private autorefreshService: AutorefreshService,
              private notificationService: NotificationService,
              private router: Router,
              private activatedRoute: ActivatedRoute,
              private titleService: Title,
              private patientService: ExternalApiService) {

    autorefreshService.start();
    this.notificationService.logToConsole();
    this.notificationService.logToSnackBar();

    const routes = this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map(() => this.activatedRoute),
      mergeMap(route => {
        // Go to second level as first level is always /
        route = route.firstChild;
        return route?.children || [];
      }));

    this.title = routes.pipe(
      filter(route => route.outlet === 'primary'),
      mergeMap(route => route.data),
      map(data => data.title)
    );
    this.title.subscribe(title => this.titleService.setTitle(`${title} - Document Manager`));

    this.navColor = routes.pipe(
      filter(route => route.outlet === 'primary'),
      mergeMap(route => route.data),
      map(data => data.color)
    );
  }

  onSelectPatient(event: Patient) {
    if (event) {
      this.navigateToPatient(event);
    }
  }

  navigateToCurrentPatient() {
    this.patientService.getSelectedPatient().subscribe(patient => this.navigateToPatient(patient));
  }

  navigateToPatient(patient: Patient) {
    this.router.navigate(['/patient', patient.id]);
  }
}
