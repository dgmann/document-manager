import {RouterModule, Routes} from '@angular/router';
import {NavigationComponent} from './navigation/navigation.component';
import {PatientComponent} from './patient.component';

const PHYSICIAN_ROUTES: Routes = [
  {
    path: ':id',
    component: PatientComponent,
    children: [],
    data: {title: 'Patient', color: 'purple'}
  },
  {
    path: '',
    component: NavigationComponent,
    outlet: 'navigation'
  }
];

export const PatientRouterModule = RouterModule.forChild(PHYSICIAN_ROUTES);
