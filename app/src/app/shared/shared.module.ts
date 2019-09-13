import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatChipsModule} from '@angular/material/chips';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatDialogModule} from '@angular/material/dialog';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {MatMenuModule} from '@angular/material/menu';
import {MatSelectModule} from '@angular/material/select';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatMomentDateModule} from '@angular/material-moment-adapter';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatSortModule} from '@angular/material/sort';
import {MatTableModule} from '@angular/material/table';
import {ActionBarComponent} from './action-bar/action-bar.component';
import {ActionMenuComponent} from './action-menu/action-menu.component';
import {AutocompleteChipsComponent} from './autocomplete-chips/autocomplete-chips.component';
import {CommentDialogComponent} from './comment-dialog/comment-dialog.component';
import {DocumentEditDialogComponent} from './document-edit-dialog/document-edit-dialog.component';
import {DocumentListComponent} from './document-list/document-list.component';
import {PatientSearchComponent} from './patient-search/patient-search.component';
import {RecordViewerComponent} from './record-viewer/record-viewer.component';
import {SplitPanelComponent} from './split-panel/split-panel.component';
import {HttpClientModule} from '@angular/common/http';
import {DragDropModule} from '@angular/cdk/drag-drop';
import {NgSelectModule} from '@ng-select/ng-select';
import {MessageBoxComponent} from './message-box/message-box.component';


@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
    FormsModule,
    NgSelectModule,
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
