import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PatientComponent} from './patient.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {RecordService} from '@app/core/records';
import {PatientService} from './patient.service';
import {CategoryService} from '@app/core/categories';
import {RouterTestingModule} from '@angular/router/testing';
import {of} from 'rxjs';
import createSpy = jasmine.createSpy;

describe('PatientComponent', () => {
  let component: PatientComponent;
  let fixture: ComponentFixture<PatientComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule
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
