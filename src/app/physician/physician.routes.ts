import {RouterModule, Routes} from "@angular/router";
import {NavigationComponent} from "./navigation/navigation.component";
import {PhysicianComponent} from "./physician.component";

const PHYSICIAN_ROUTES: Routes = [
  {
    path: '',
    component: PhysicianComponent
  },
  {
    path: '',
    component: NavigationComponent,
    outlet: 'navigation'
  },
];

export const PhysicianRouterModule = RouterModule.forChild(PHYSICIAN_ROUTES);
