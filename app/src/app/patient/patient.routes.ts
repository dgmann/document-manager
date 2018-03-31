import {RouterModule, Routes} from "@angular/router";
import {NavigationComponent} from "./navigation/navigation.component";
import {PatientComponent} from "./patient.component";

const PHYSICIAN_ROUTES: Routes = [
  {
    path: '',
    component: PatientComponent,
    children: [],
    data: {title: 'Patient'}
  },
  {
    path: '',
    component: NavigationComponent,
    outlet: 'navigation'
  }
];

export const PatientRouterModule = RouterModule.forChild(PHYSICIAN_ROUTES);
