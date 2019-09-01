import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageOtherComponent} from './page-other.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {PhysicianService} from '../physician.service';
import {of} from 'rxjs';

describe('PageOtherComponent', () => {
  let component: PageOtherComponent;
  let fixture: ComponentFixture<PageOtherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        PageOtherComponent
      ],
      providers: [{provide: PhysicianService, useValue: {otherRecords$: of(), selectedIds$: of()}}],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PageOtherComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
