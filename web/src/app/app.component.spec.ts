import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import {AppComponent} from './app.component';
import {AppModule} from './app.module';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {APP_BASE_HREF} from '@angular/common';
import { NotificationService } from './core/notifications';
import { AutorefreshService } from './core/autorefresh';
import { ExternalApiService } from './shared/document-edit-dialog/external-api.service';
import { of } from 'rxjs';

describe('AppComponent', () => {
  let component: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule
      ],
      providers: [
        {provide: APP_BASE_HREF, useValue: '/'},
        {provide: NotificationService, useValue: {
          logToConsole: jest.fn(),
          logToSnackBar: jest.fn()
        }},
        {provide: AutorefreshService, useValue: {
          start: jest.fn()
        }},
        {provide: ExternalApiService, useValue: {
          getSelectedPatient: jest.fn().mockReturnValue(of())
        }}
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
