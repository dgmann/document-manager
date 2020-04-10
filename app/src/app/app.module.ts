import {BrowserModule} from '@angular/platform-browser';
import localeDe from '@angular/common/locales/de';
import {APP_INITIALIZER, LOCALE_ID, NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {SharedModule} from '@app/shared';
import {CoreModule} from '@app/core/core.module';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {registerLocaleData} from '@angular/common';
import {AppRoutingModule} from './app-routing.module';
import {MatMenuModule} from '@angular/material/menu';
import {ConfigService} from '@app/core/config';

registerLocaleData(localeDe, 'de');

const initializerConfigFn = (appConfig: ConfigService) => {
  return () => {
    return appConfig.loadAppConfig();
  };
};

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    SharedModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    MatIconModule,
    MatButtonModule,
    MatSidenavModule,
    MatToolbarModule,
    MatSnackBarModule,
    CoreModule,
    MatMenuModule,
  ],
  providers: [
    {provide: LOCALE_ID, useValue: 'de-DE'},
    {
      provide: APP_INITIALIZER,
      useFactory: initializerConfigFn,
      multi: true,
      deps: [ConfigService],
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
