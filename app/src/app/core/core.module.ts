import {NgModule, Optional, SkipSelf} from '@angular/core';
import {CommonModule} from '@angular/common';
import {throwIfAlreadyLoaded} from './module-import-guard';
import {NgDragDropModule} from 'ng-drag-drop';
import {StoreModule} from './store';
import {EventSnackbarComponent} from './event-snackbar/event-snackbar.component';

@NgModule({
  imports: [
    CommonModule,
    StoreModule,
    NgDragDropModule.forRoot(),
  ],
  declarations: [EventSnackbarComponent],
  entryComponents: [
    EventSnackbarComponent
  ],
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    throwIfAlreadyLoaded(parentModule, 'CoreModule');
  }
}
