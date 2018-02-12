import {registerLocaleData} from '@angular/common';
import localeDe from '@angular/common/locales/de';
import {LOCALE_ID, NgModule} from '@angular/core';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FormsModule} from "@angular/forms";
import {
  MatAutocompleteModule,
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
import {NgDragDropModule} from "ng-drag-drop";


import {AppComponent} from './app.component';
import {DocumentEditDialogComponent} from "./document-edit-dialog/document-edit-dialog.component";
import {InboxModule} from "./inbox/index";
import {PatientService} from "./shared";
import {NotificationService} from "./shared/notification-service";
import {TagService} from "./shared/tag-service";
import {WebsocketService} from "./shared/websocket-service";
import {StoreModule} from "./store";


registerLocaleData(localeDe, 'de');

@NgModule({
  declarations: [
    AppComponent,
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
    MatAutocompleteModule,
    StoreModule,
    NgDragDropModule.forRoot(),
    InboxModule
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
    PatientService,
    WebsocketService,
    NotificationService,
    TagService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
