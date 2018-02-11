import {registerLocaleData} from '@angular/common';
import localeDe from '@angular/common/locales/de';
import {LOCALE_ID, NgModule} from '@angular/core';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FormsModule} from "@angular/forms";
import {
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatSnackBarModule
} from "@angular/material";
import {MatMomentDateModule} from '@angular/material-moment-adapter';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatToolbarModule} from '@angular/material/toolbar';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {DndModule} from "ng2-dnd";


import {AppComponent} from './app.component';
import {DocumentEditDialogComponent} from "./document-edit-dialog/document-edit-dialog.component";
import {DocumentListModule} from "./document-list";
import {RecordViewerComponent} from "./record-viewer/record-viewer.component";
import {PatientService} from "./shared";
import {NotificationService} from "./shared/notification-service";
import {WebsocketService} from "./shared/websocket-service";
import {StoreModule} from "./store";


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
    StoreModule,
    DndModule.forRoot()
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
    PatientService,
    WebsocketService,
    NotificationService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
