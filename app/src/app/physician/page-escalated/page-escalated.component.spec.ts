import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageEscalatedComponent} from './page-escalated.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {PhysicianService} from '../physician.service';
import {of} from 'rxjs';

describe('PageEscalatedComponent', () => {
  let component: PageEscalatedComponent;
  let fixture: ComponentFixture<PageEscalatedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        PageEscalatedComponent
      ],
      providers: [{provide: PhysicianService, useValue: {escalatedRecords$: of([])}}],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PageEscalatedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
