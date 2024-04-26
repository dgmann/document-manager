import {NgModule} from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { MatMenuModule } from '@angular/material/menu';
import {EditorComponent} from './editor.component';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {EditorRouterModule} from './editor.routes';
import {SharedModule} from '../shared';
import {EditorGuard} from './editor.guard';
import {PageListComponent} from './page-list/page-list.component';
import {DragDropModule} from '@angular/cdk/drag-drop';

@NgModule({
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatMenuModule,
    MatIconModule,
    SharedModule,
    EditorRouterModule,
    DragDropModule,
    NgOptimizedImage
  ],
  declarations: [EditorComponent, PageListComponent],
  exports: [EditorComponent],
  providers: [
    EditorGuard
  ]
})
export class EditorModule {
}
