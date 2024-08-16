import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import {MAT_DIALOG_DATA, MatDialogModule} from '@angular/material/dialog';

import { MessageBoxComponent } from './message-box.component';

describe('MessageBoxComponent', () => {
  let component: MessageBoxComponent;
  let fixture: ComponentFixture<MessageBoxComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [MatDialogModule],
      providers: [{ provide: MAT_DIALOG_DATA, useValue: { title: 'Test', text: 'Test Text'} }],
      declarations: [MessageBoxComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageBoxComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
