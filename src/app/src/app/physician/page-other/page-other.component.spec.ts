import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageOtherComponent} from './page-other.component';

describe('PageOtherComponent', () => {
  let component: PageOtherComponent;
  let fixture: ComponentFixture<PageOtherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PageOtherComponent]
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
