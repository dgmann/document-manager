import {RouterModule, Routes} from "@angular/router";
import {InboxComponent} from "./inbox.component";

const INBOX_ROUTES: Routes = [
  {
    path: '',
    component: InboxComponent,
    data: {title: 'Inbox'}
  },
];

export const InboxRouterModule = RouterModule.forChild(INBOX_ROUTES);
