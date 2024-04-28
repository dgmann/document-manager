import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import {MatBottomSheetModule} from '@angular/material/bottom-sheet';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import { provideMockStore } from '@ngrx/store/testing';

import {InboxComponent} from './inbox.component';
import {InboxService} from './inbox.service';
import {of} from 'rxjs';
import {SharedModule} from '../shared';
import {RecordService} from '../core/records';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {NoopAnimationsModule} from '@angular/platform-browser/animations';
import { RouterModule } from '@angular/router';

describe('InboxComponent', () => {
  let component: InboxComponent;
  let fixture: ComponentFixture<InboxComponent>;
  let inboxService;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule,
        RouterModule,
        NoopAnimationsModule,
        MatBottomSheetModule,
        MatSnackBarModule,
      ],
      declarations: [InboxComponent],
      providers: [
        provideMockStore(),
        {provide: InboxService, useValue: {
          allInboxRecords$: of([]),
          selectedIds$: of([]),
          selectedRecords$: of([]),
          isMultiSelect$: of(false),
          loadRecords: jest.fn(),
          upload: jest.fn(),
          selectIds: jest.fn(),
          deleteSelectedRecords: jest.fn(),
          updateSelectedRecords: jest.fn()
        },
      },
        {provide: RecordService, useValue: {
          updatePages: jest.fn()
        }}
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    }).compileComponents();
    inboxService = TestBed.get(InboxService);
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InboxComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
    expect(inboxService.loadRecords).toHaveBeenCalled();
  });

  it('should upload pdfs on drop', () => {
    const files = [
      {},
      {},
      {}
    ];
    const event = {
      dataTransfer: {
        files,
        clearData: () => {}
      },
      preventDefault: () => {},
      stopPropagation: () => {}
    } as unknown as DragEvent;

    component.onDrop(event);
    expect(inboxService.upload).toHaveBeenCalledTimes(3);
  });
});
