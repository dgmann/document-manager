import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {TagListComponent} from './tag-list.component';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {MatTableModule} from '@angular/material/table';
import {of} from 'rxjs';

describe('TagListComponent', () => {
  let component: TagListComponent;
  let fixture: ComponentFixture<TagListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [MatTableModule],
      declarations: [TagListComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TagListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
