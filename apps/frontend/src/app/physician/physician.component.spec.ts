import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {PhysicianComponent} from './physician.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {PhysicianService} from './physician.service';
import {of} from 'rxjs';
import { RouterModule } from '@angular/router';

describe('PhysicianComponent', () => {
  let component: PhysicianComponent;
  let fixture: ComponentFixture<PhysicianComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterModule.forRoot([])
      ],
      declarations: [
        PhysicianComponent
      ],
      providers: [{
        provide: PhysicianService, useValue: {
          selectedRecords$: of(),
          selectedIds$: of(),
          load: jest.fn(),
          selectIds: jest.fn()
        }
      }
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhysicianComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
