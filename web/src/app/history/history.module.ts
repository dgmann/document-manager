import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MatButtonModule} from '@angular/material/button';
import {HistoryComponent} from './history.component';
import {SharedModule} from '@app/shared';
import {HistoryRouterModule} from '@app/history/history.routes';
import {HistoryService} from '@app/history/history-service';


@NgModule({
  declarations: [HistoryComponent],
    imports: [
        CommonModule,
        HistoryRouterModule,
        SharedModule,
        MatButtonModule
    ],
  providers: [
    HistoryService
  ]
})
export class HistoryModule {
}
