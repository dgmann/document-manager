import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RecordFilterComponent } from './record-filter.component';

describe('RecordFilterComponent', () => {
  let component: RecordFilterComponent;
  let fixture: ComponentFixture<RecordFilterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [RecordFilterComponent]
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
