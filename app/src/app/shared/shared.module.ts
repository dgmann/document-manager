import {CommonModule} from "@angular/common";
import {NgModule} from "@angular/core";
import {FlexLayoutModule} from "@angular/flex-layout";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {
  MatAutocompleteModule,
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatMenuModule,
  MatSelectModule,
  MatTooltipModule
} from "@angular/material";
import {MatMomentDateModule} from "@angular/material-moment-adapter";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatSortModule} from "@angular/material/sort";
import {MatTableModule} from "@angular/material/table";
import {NgDragDropModule} from "ng-drag-drop";
import {NgxMatSelectSearchModule} from "ngx-mat-select-search";
import {ActionMenuComponent} from './action-menu/action-menu.component';
import {AutocompleteChipsComponent} from "./autocomplete-chips/autocomplete-chips.component";
import {AutocompleteInputComponent} from "./autocomplete-input/autocomplete-input.component";
import {CategoryService} from "./category-service";
import {DocumentEditDialogComponent} from "./document-edit-dialog/document-edit-dialog.component";
import {DocumentListComponent} from "./document-list/document-list.component";
import {ExternalApiService} from "./external-api.service";
import {NotificationService} from "./notification-service";
import {PatientSearchComponent} from './patient-search/patient-search.component';
import {RecordViewerComponent} from "./record-viewer/record-viewer.component";
import {SplitPanelComponent} from './split-panel/split-panel.component';
import {TagService} from "./tag-service";
import {WebsocketService} from "./websocket-service";


@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    FlexLayoutModule,
    MatSelectModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule,
    MatSortModule,
    MatTableModule,
    MatTooltipModule,
    NgDragDropModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    MatInputModule,
    MatMomentDateModule,
    MatDatepickerModule,
    MatChipsModule,
    MatMenuModule,
    NgxMatSelectSearchModule
  ],
  declarations: [
    RecordViewerComponent,
    DocumentEditDialogComponent,
    DocumentListComponent,
    PatientSearchComponent,
    SplitPanelComponent,
    AutocompleteInputComponent,
    AutocompleteChipsComponent,
    ActionMenuComponent
  ],
  entryComponents: [
    DocumentEditDialogComponent
  ],
  exports: [
    CommonModule,
    FlexLayoutModule,
    DocumentEditDialogComponent,
    RecordViewerComponent,
    DocumentListComponent,
    PatientSearchComponent,
    SplitPanelComponent
  ],
  providers: []
})
export class SharedModule {
  static forRoot() {
    return {
      ngModule: SharedModule,
      providers: [
        ExternalApiService,
        WebsocketService,
        NotificationService,
        TagService,
        CategoryService
      ]
    }
  }
}
