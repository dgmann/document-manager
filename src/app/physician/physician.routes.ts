import {RouterModule, Routes} from "@angular/router";
import {NavigationComponent} from "./navigation/navigation.component";
import {PageEscalatedComponent} from "./page-escalated/page-escalated.component";
import {PageOtherComponent} from "./page-other/page-other.component";
import {PageReviewComponent} from "./page-review/page-review.component";
import {PhysicianComponent} from "./physician.component";

const PHYSICIAN_ROUTES: Routes = [
  {
    path: '',
    component: PhysicianComponent,
    children: [
      {
        path: 'review',
        component: PageReviewComponent,
      },
      {
        path: 'escalated',
        component: PageEscalatedComponent,
      },
      {
        path: 'other',
        component: PageOtherComponent,
      }
    ]
  },
  {
    path: '',
    component: NavigationComponent,
    outlet: 'navigation'
  }
];

export const PhysicianRouterModule = RouterModule.forChild(PHYSICIAN_ROUTES);
