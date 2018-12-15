import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { FlexLayoutModule } from "@angular/flex-layout";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
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
import { MatMomentDateModule } from "@angular/material-moment-adapter";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatSortModule } from "@angular/material/sort";
import { MatTableModule } from "@angular/material/table";
import { NgDragDropModule } from "ng-drag-drop";
import { ActionBarComponent } from './action-bar/action-bar.component';
import { ActionMenuComponent } from './action-menu/action-menu.component';
import { AutocompleteChipsComponent } from "./autocomplete-chips/autocomplete-chips.component";
import { AutocompleteInputComponent } from "./autocomplete-input/autocomplete-input.component";
import { CommentDialogComponent } from './comment-dialog/comment-dialog.component';
import { DocumentEditDialogComponent } from "./document-edit-dialog/document-edit-dialog.component";
import { DocumentListComponent } from "./document-list/document-list.component";
import { PatientSearchComponent } from './patient-search/patient-search.component';
import { RecordViewerComponent } from "./record-viewer/record-viewer.component";
import { SplitPanelComponent } from './split-panel/split-panel.component';
import { HttpClientModule } from "@angular/common/http";
import { DragDropModule } from "@angular/cdk/drag-drop";
import { MessageBoxComponent } from './message-box/message-box.component';


@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
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
    DragDropModule
  ],
  declarations: [
    RecordViewerComponent,
    DocumentEditDialogComponent,
    DocumentListComponent,
    PatientSearchComponent,
    SplitPanelComponent,
    AutocompleteInputComponent,
    AutocompleteChipsComponent,
    ActionMenuComponent,
    ActionBarComponent,
    CommentDialogComponent,
    MessageBoxComponent
  ],
  entryComponents: [
    DocumentEditDialogComponent,
    CommentDialogComponent,
    MessageBoxComponent
  ],
  exports: [
    CommonModule,
    FlexLayoutModule,
    DocumentEditDialogComponent,
    RecordViewerComponent,
    DocumentListComponent,
    PatientSearchComponent,
    SplitPanelComponent,
    ActionBarComponent
  ],
  providers: []
})
export class SharedModule {
}
