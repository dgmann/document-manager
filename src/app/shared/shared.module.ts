import {CommonModule} from "@angular/common";
import {NgModule} from "@angular/core";
import {FlexLayoutModule} from "@angular/flex-layout";
import {FormsModule} from "@angular/forms";
import {
  MatAutocompleteModule,
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatMenuModule,
  MatTooltipModule
} from "@angular/material";
import {MatMomentDateModule} from "@angular/material-moment-adapter";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatSortModule} from "@angular/material/sort";
import {MatTableModule} from "@angular/material/table";
import {NgDragDropModule} from "ng-drag-drop";
import {DocumentEditDialogComponent} from "./document-edit-dialog/document-edit-dialog.component";
import {DocumentListComponent} from "./document-list/document-list.component";
import {NotificationService} from "./notification-service";
import {PatientService} from "./patient-service";
import {RecordViewerComponent} from "./record-viewer/record-viewer.component";
import {TagService} from "./tag-service";
import {WebsocketService} from "./websocket-service";


@NgModule({
  imports: [
    CommonModule,
    FormsModule,
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
    MatFormFieldModule,
    MatInputModule,
    MatMomentDateModule,
    MatDatepickerModule,
    MatChipsModule,
    MatMenuModule
  ],
  declarations: [
    RecordViewerComponent,
    DocumentEditDialogComponent,
    DocumentListComponent
  ],
  entryComponents: [
    DocumentEditDialogComponent
  ],
  exports: [
    CommonModule,
    FlexLayoutModule,
    DocumentEditDialogComponent,
    RecordViewerComponent,
    DocumentListComponent
  ],
  providers: []
})
export class SharedModule {
  static forRoot() {
    return {
      ngModule: SharedModule,
      providers: [
        PatientService,
        WebsocketService,
        NotificationService,
        TagService
      ]
    }
  }
}
