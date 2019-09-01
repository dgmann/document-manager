import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';


const routes: Routes = [
  {
    path: 'inbox',
    loadChildren: () => import('@app/inbox/inbox.module').then(mod => mod.InboxModule)
  },
  {
    path: 'physician',
    loadChildren: () => import('@app/physician/physician.module').then(mod => mod.PhysicianModule)
  },
  {
    path: 'patient',
    loadChildren: () => import('@app/patient/patient.module').then(mod => mod.PatientModule)
  },
  {
    path: 'editor',
    loadChildren: () => import('@app/editor/editor.module').then(mod => mod.EditorModule)
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
