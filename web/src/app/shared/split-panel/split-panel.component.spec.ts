import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {SplitPanelComponent} from './split-panel.component';

describe('SplitPanelComponent', () => {
  let component: SplitPanelComponent;
  let fixture: ComponentFixture<SplitPanelComponent>;

  beforeEach(waitForAsync(() => {
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
