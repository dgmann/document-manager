import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageReviewComponent} from './page-review.component';

describe('PageReviewComponent', () => {
  let component: PageReviewComponent;
  let fixture: ComponentFixture<PageReviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PageReviewComponent]
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
