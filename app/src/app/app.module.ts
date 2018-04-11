import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';
import { LOCALE_ID, NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatButtonModule, MatIconModule, MatSnackBarModule } from "@angular/material";
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgDragDropModule } from "ng-drag-drop";


import { AppComponent } from './app.component';
import { AppRoutesModule } from "./app.router";
import { SharedModule } from "./shared";
import { StoreModule } from "./store";
import { DndModule } from "ng2-dnd";


registerLocaleData(localeDe, 'de');

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutesModule,
    FlexLayoutModule,
    MatIconModule,
    MatButtonModule,
    MatSidenavModule,
    MatToolbarModule,
    MatSnackBarModule,
    StoreModule,
    NgDragDropModule.forRoot(),
    SharedModule.forRoot(),
    DndModule.forRoot()
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
