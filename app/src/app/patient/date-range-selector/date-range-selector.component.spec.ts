import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DateRangeSelectorComponent } from './date-range-selector.component';
import { FormsModule } from "@angular/forms";
import {
  MatDatepickerModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatNativeDateModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";

describe('DateRangeSelectorComponent', () => {
  let component: DateRangeSelectorComponent;
  let fixture: ComponentFixture<DateRangeSelectorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        BrowserAnimationsModule,
        FormsModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule
      ],
      declarations: [DateRangeSelectorComponent]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DateRangeSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
