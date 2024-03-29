import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {RecordFilterComponent} from './record-filter.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {of} from 'rxjs';

describe('RecordFilterComponent', () => {
  let component: RecordFilterComponent;
  let fixture: ComponentFixture<RecordFilterComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [
        RecordFilterComponent
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RecordFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
