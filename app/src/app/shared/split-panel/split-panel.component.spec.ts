import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {SplitPanelComponent} from './split-panel.component';

describe('SplitPanelComponent', () => {
  let component: SplitPanelComponent;
  let fixture: ComponentFixture<SplitPanelComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [SplitPanelComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SplitPanelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
