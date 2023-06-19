import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import {NoopAnimationsModule} from '@angular/platform-browser/animations';
import {ConfigService} from '@app/core/config';

import {PatientComponent} from './patient.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {RecordService} from '@app/core/records';
import {PatientService} from './patient.service';
import {CategoryService} from '@app/core/categories';
import {RouterTestingModule} from '@angular/router/testing';
import {of} from 'rxjs';
import createSpy = jasmine.createSpy;
import {SharedModule} from '@app/shared';

describe('PatientComponent', () => {
  let component: PatientComponent;
  let fixture: ComponentFixture<PatientComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule,
        HttpClientTestingModule,
        SharedModule,
        NoopAnimationsModule
      ],
      declarations: [PatientComponent],
      providers: [
        {provide: RecordService, useValue: {}},
        {
          provide: PatientService, useValue: {
            selectedCategory$: of(),
            filteredPatientRecord$: of(),
            selectedPatient$: of(),
            selectedRecord$: of(),
            selectedId$: of(),
            selectPatient: createSpy(),
            selectCategory: createSpy()
          }
        },
        {
          provide: CategoryService, useValue: {
            categoryMap: of(),
            load: createSpy()
          }
        },
        { provide: ConfigService, useValue: {getApiUrl: () => 'http://test.com'}}
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PatientComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
