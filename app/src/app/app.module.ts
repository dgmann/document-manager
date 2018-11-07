import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';
import { ErrorHandler, LOCALE_ID, NgModule } from '@angular/core';
import { MatButtonModule, MatIconModule, MatSnackBarModule } from "@angular/material";
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceWorkerModule } from '@angular/service-worker'
import BugsnagErrorHandler from 'bugsnag-angular'
import bugsnag from 'bugsnag-js';
import { environment } from "../environments/environment";


import { AppComponent } from './app.component';
import { AppRoutesModule } from "./app.router";
import { CoreModule } from "./core/core.module";
import { SharedModule } from "./shared";


const bugsnagClient = bugsnag(environment.bugsnagKey);

export function errorHandlerFactory() {
  return new BugsnagErrorHandler(bugsnagClient)
}
registerLocaleData(localeDe, 'de');

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    SharedModule,
    BrowserAnimationsModule,
    AppRoutesModule,
    MatIconModule,
    MatButtonModule,
    MatSidenavModule,
    MatToolbarModule,
    MatSnackBarModule,
    CoreModule,
    ServiceWorkerModule.register('/ngsw-worker.js', {enabled: environment.production})
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
    {provide: ErrorHandler, useFactory: errorHandlerFactory}
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
