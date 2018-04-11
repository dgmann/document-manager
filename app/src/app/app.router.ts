import { RouterModule, Routes } from "@angular/router";

export const ROUTE_CONFIG: Routes = [
  {
    path: 'inbox',
    loadChildren: 'app/inbox/inbox.module#InboxModule'
  },
  {
    path: 'physician',
    loadChildren: 'app/physician/physician.module#PhysicianModule'
  },
  {
    path: 'patient',
    loadChildren: 'app/patient/patient.module#PatientModule'
  },
  {
    path: 'editor',
    loadChildren: 'app/editor/editor.module#EditorModule'
  },
  {
    path: '',
    redirectTo: 'inbox',
    pathMatch: 'full'
  }
];

export const AppRoutesModule = RouterModule.forRoot(ROUTE_CONFIG);
