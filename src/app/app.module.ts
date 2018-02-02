import { BrowserModule } from '@angular/platform-browser';
import { LOCALE_ID, NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatSidenavModule } from '@angular/material/sidenav';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatToolbarModule } from '@angular/material/toolbar';
import {
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatSnackBarModule
} from "@angular/material";
import { MatMomentDateModule } from '@angular/material-moment-adapter';
import { MatButtonModule } from '@angular/material/button';
import { FormsModule } from "@angular/forms";
import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';


import { AppComponent } from './app.component';
import { DocumentListModule } from "./document-list";
import { RecordViewerComponent } from "./record-viewer/record-viewer.component";
import { DocumentEditDialogComponent } from "./document-edit-dialog/document-edit-dialog.component";
import { PatientService } from "./shared";
import { StoreModule } from "./store";


registerLocaleData(localeDe, 'de');

@NgModule({
  declarations: [
    AppComponent,
    RecordViewerComponent,
    DocumentEditDialogComponent
  ],
  entryComponents: [
    DocumentEditDialogComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    FlexLayoutModule,
    MatCardModule,
    MatSidenavModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatMomentDateModule,
    MatDatepickerModule,
    MatChipsModule,
    MatSnackBarModule,
    DocumentListModule,
    StoreModule
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
    PatientService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
