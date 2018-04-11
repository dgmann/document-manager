import { RouterModule, Routes } from "@angular/router";
import { EditorComponent } from "./editor.component";
import { EditorGuard } from "./editor.guard";

const INBOX_ROUTES: Routes = [
  {
    path: ':id',
    canActivate: [EditorGuard],
    component: EditorComponent,
    data: {title: 'Editor'}
  }
];

export const EditorRouterModule = RouterModule.forChild(INBOX_ROUTES);
