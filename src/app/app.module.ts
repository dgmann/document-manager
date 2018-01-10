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
  MatNativeDateModule
} from "@angular/material";
import { MatButtonModule } from '@angular/material/button';
import { StoreModule } from "@ngrx/store";
import { StoreDevtoolsModule } from "@ngrx/store-devtools";
import { FormsModule } from "@angular/forms";
import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';


import { AppComponent } from './app.component';
import { DocumentListModule } from "./document-list";
import { metaReducers, reducers } from './reducers';
import { ApiModule } from "./api";
import { EffectsModule } from "@ngrx/effects";
import { RecordViewerComponent } from "./record-viewer/record-viewer.component";
import { DocumentEditDialogComponent } from "./document-edit-dialog/document-edit-dialog.component";


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
    MatNativeDateModule,
    MatDatepickerModule,
    MatChipsModule,
    DocumentListModule,
    StoreModule.forRoot(reducers, { metaReducers }),
    StoreDevtoolsModule.instrument({
      maxAge: 25 //  Retains last 25 states
    }),
    EffectsModule.forRoot([]),
    ApiModule
  ],
  providers: [{provide: LOCALE_ID, useValue: 'de-DE'}],
  bootstrap: [AppComponent]
})
export class AppModule {}
