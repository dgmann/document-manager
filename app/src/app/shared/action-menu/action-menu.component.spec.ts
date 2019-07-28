import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionMenuComponent } from './action-menu.component';
import { MatMenuModule } from "@angular/material/menu";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { Status } from "../../core/store";

describe('ActionMenuComponent', () => {
  let component: ActionMenuComponent;
  let fixture: ComponentFixture<ActionMenuComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [MatMenuModule],
      declarations: [ActionMenuComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionMenuComponent);
    component = fixture.componentInstance;
    component.record = {
      status: Status.ESCALATED
    } as any;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
