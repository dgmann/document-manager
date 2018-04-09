import {RouterModule, Routes} from "@angular/router";
import {InboxComponent} from "./inbox.component";
import {NavigationComponent} from "./navigation/navigation.component";

const INBOX_ROUTES: Routes = [
  {
    path: '',
    component: InboxComponent,
    data: {title: 'Inbox'}
  },
  {
    path: '',
    component: NavigationComponent,
    outlet: 'navigation'
  }
];

export const InboxRouterModule = RouterModule.forChild(INBOX_ROUTES);
