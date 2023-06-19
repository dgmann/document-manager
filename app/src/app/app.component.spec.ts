import { TestBed, waitForAsync } from '@angular/core/testing';
import {AppComponent} from './app.component';
import {AppModule} from './app.module';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {APP_BASE_HREF} from '@angular/common';

describe('AppComponent', () => {
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule
      ],
      providers: [{provide: APP_BASE_HREF, useValue: '/'}],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    }).compileComponents();
  }));
});
