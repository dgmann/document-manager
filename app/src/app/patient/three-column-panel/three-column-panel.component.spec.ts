import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ThreeColumnPanelComponent} from './three-column-panel.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {FlexLayoutModule} from '@angular/flex-layout';

describe('ThreeColumnPanelComponent', () => {
  let component: ThreeColumnPanelComponent;
  let fixture: ComponentFixture<ThreeColumnPanelComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [FlexLayoutModule],
      declarations: [ThreeColumnPanelComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
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
