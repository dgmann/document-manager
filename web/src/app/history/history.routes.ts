import {RouterModule, Routes} from '@angular/router';
import {HistoryComponent} from '@app/history/history.component';

const HISTORY_ROUTES: Routes = [
  {
    path: '',
    component: HistoryComponent,
    data: {title: 'History'}
  }
];

export const HistoryRouterModule = RouterModule.forChild(HISTORY_ROUTES);
