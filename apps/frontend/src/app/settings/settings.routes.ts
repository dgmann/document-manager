import {RouterModule, Routes} from '@angular/router';
import {SettingsComponent} from '@app/settings/settings.component';

const SETTINGS_ROUTES: Routes = [
  {
    path: '',
    component: SettingsComponent,
    data: {title: 'Einstellungen'},
  }
];

export const SettingsRouterModule = RouterModule.forChild(SETTINGS_ROUTES);
