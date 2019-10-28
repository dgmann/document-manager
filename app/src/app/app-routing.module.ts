import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';


const routes: Routes = [
  {
    path: 'inbox',
    loadChildren: () => import('@app/inbox').then(mod => mod.InboxModule)
  },
  {
    path: 'history',
    loadChildren: () => import('@app/history').then(mod => mod.HistoryModule)
  },
  {
    path: 'physician',
    loadChildren: () => import('@app/physician').then(mod => mod.PhysicianModule)
  },
  {
    path: 'patient',
    loadChildren: () => import('@app/patient').then(mod => mod.PatientModule)
  },
  {
    path: 'editor',
    loadChildren: () => import('@app/editor').then(mod => mod.EditorModule)
  },
  {
    path: 'settings',
    loadChildren: () => import('@app/settings').then(mod => mod.SettingsModule)
  },
  {
    path: '',
    redirectTo: 'inbox',
    pathMatch: 'full'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
