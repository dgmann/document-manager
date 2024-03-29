import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { MatDateFnsModule } from '@angular/material-date-fns-adapter';
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
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatSortModule} from '@angular/material/sort';
import {MatTableModule} from '@angular/material/table';
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
import {MessageBoxComponent} from './message-box/message-box.component';
import {PdfLinkPipe} from './pdf-link/pdf-link.pipe';
import {IdsPipe} from './pipes/ids/ids.pipe';
import {CategoryPipe} from './pipes/category/category.pipe';
import { PatientPipe } from './pipes/patient/patient.pipe';
import { ContainsPipe } from './pipes/contains/contains.pipe';


@NgModule({
    imports: [
        CommonModule,
        HttpClientModule,
        FormsModule,
        ReactiveFormsModule,
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
        MatDateFnsModule,
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
        CommentDialogComponent,
        MessageBoxComponent,
        PdfLinkPipe,
        IdsPipe,
        CategoryPipe,
        PatientPipe,
        ContainsPipe
    ],
    exports: [
        CommonModule,
        DocumentEditDialogComponent,
        RecordViewerComponent,
        DocumentListComponent,
        PatientSearchComponent,
        SplitPanelComponent,
        PdfLinkPipe,
        IdsPipe,
        CategoryPipe,
        PatientPipe,
    ],
    providers: []
})
export class SharedModule {
}
