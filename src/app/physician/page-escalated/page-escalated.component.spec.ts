import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {PageEscalatedComponent} from './page-escalated.component';

describe('PageEscalatedComponent', () => {
  let component: PageEscalatedComponent;
  let fixture: ComponentFixture<PageEscalatedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PageEscalatedComponent]
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
