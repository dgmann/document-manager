import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ThreeColumnPanelComponent } from './three-column-panel.component';

describe('ThreeColumnPanelComponent', () => {
  let component: ThreeColumnPanelComponent;
  let fixture: ComponentFixture<ThreeColumnPanelComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ThreeColumnPanelComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ThreeColumnPanelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
