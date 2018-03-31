import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {MultiRecordListComponent} from './multi-record-list.component';

describe('MultiRecordListComponent', () => {
  let component: MultiRecordListComponent;
  let fixture: ComponentFixture<MultiRecordListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [MultiRecordListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiRecordListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
