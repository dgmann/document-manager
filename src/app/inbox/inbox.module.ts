import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {FlexLayoutModule} from "@angular/flex-layout";
import {MatAutocompleteModule, MatDialogModule, MatIconModule, MatTooltipModule,} from "@angular/material";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatSortModule} from "@angular/material/sort";
import {MatTableModule} from "@angular/material/table";
import {EffectsModule} from "@ngrx/effects";
import {StoreModule} from "@ngrx/store";
import {NgDragDropModule} from "ng-drag-drop";
import {DocumentListComponent} from "./document-list/document-list.component";
import {InboxComponent} from './inbox.component';
import {InboxService} from "./inbox.service";
import {RecordViewerComponent} from "./record-viewer/record-viewer.component";
import {metaReducers, reducers} from './reducers';
import {InboxEffects} from "./store/inbox.effects";

@NgModule({
  imports: [
    CommonModule,
    FlexLayoutModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule,
    MatSortModule,
    MatTableModule,
    MatTooltipModule,
    NgDragDropModule,
    MatAutocompleteModule,
    StoreModule.forFeature("inbox", reducers, {metaReducers}),
    EffectsModule.forFeature([InboxEffects])
  ],
  declarations: [
    InboxComponent,
    DocumentListComponent,
    RecordViewerComponent
  ],
  exports: [
    InboxComponent
  ],
  providers: [
    InboxService
  ]
})
export class InboxModule {
}
