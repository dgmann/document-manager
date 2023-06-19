import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatDialogModule } from '@angular/material/dialog';

import {MultiRecordListComponent} from './multi-record-list.component';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {DocumentEditDialogService} from '@app/shared';
import createSpyObj = jasmine.createSpyObj;

describe('MultiRecordListComponent', () => {
  let component: MultiRecordListComponent;
  let fixture: ComponentFixture<MultiRecordListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        MatCardModule,
        MatIconModule,
        MatDialogModule
      ],
      declarations: [MultiRecordListComponent],
      providers: [
        {provide: DocumentEditDialogService, useValue: createSpyObj(['open'])}
      ]
    })
      .compileComponents();

  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiRecordListComponent);
    component = fixture.componentInstance;
    component.records = [];
    component.categories = {};
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
