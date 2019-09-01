import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageReviewComponent} from './page-review.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {MatTableModule} from '@angular/material/table';
import {PhysicianService} from '../physician.service';
import {of} from 'rxjs';

describe('PageReviewComponent', () => {
  let component: PageReviewComponent;
  let fixture: ComponentFixture<PageReviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [MatTableModule],
      declarations: [
        PageReviewComponent
      ],
      providers: [{
        provide: PhysicianService,
        useValue: {reviewRecords$: of(), selectedIds$: of()}
      }],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PageReviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
