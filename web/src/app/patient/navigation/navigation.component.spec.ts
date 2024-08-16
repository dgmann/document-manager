import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {NavigationComponent} from './navigation.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {PatientService} from '../patient.service';

describe('Patient NavigationComponent', () => {
  let component: NavigationComponent;
  let fixture: ComponentFixture<NavigationComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [NavigationComponent],
      providers: [{provide: PatientService, useValue: {getSelectedPatient: jest.fn()}}],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NavigationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
