import { HttpClientTestingModule } from '@angular/common/http/testing';
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import { provideMockConfigService } from '@app/core/config';
import { PdfLinkPipe } from '@app/shared/pdf-link/pdf-link.pipe';

import {ActionMenuComponent} from './action-menu.component';
import {MatMenuModule} from '@angular/material/menu';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {Status} from '@app/core/records';

describe('ActionMenuComponent', () => {
  let component: ActionMenuComponent;
  let fixture: ComponentFixture<ActionMenuComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [MatMenuModule],
      declarations: [ActionMenuComponent, PdfLinkPipe],
      providers: [provideMockConfigService()],
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
