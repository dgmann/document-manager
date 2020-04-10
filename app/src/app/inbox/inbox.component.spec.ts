import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {InboxComponent} from './inbox.component';
import {InboxService} from './inbox.service';
import {of} from 'rxjs';
import {SharedModule} from '../shared';
import {RecordService} from '../core/records';
import {RouterTestingModule} from '@angular/router/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {NoopAnimationsModule} from '@angular/platform-browser/animations';
import createSpy = jasmine.createSpy;
import createSpyObj = jasmine.createSpyObj;

describe('InboxComponent', () => {
  let component: InboxComponent;
  let fixture: ComponentFixture<InboxComponent>;
  let inboxService;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule,
        RouterTestingModule,
        NoopAnimationsModule
      ],
      declarations: [InboxComponent],
      providers: [{
        provide: InboxService, useValue: {
          allInboxRecords$: of([]),
          selectedIds$: of([]),
          selectedRecords$: of([]),
          isMultiSelect$: of(false),
          loadRecords: createSpy('loadRecords'),
          upload: createSpy('upload'),
          selectIds: createSpy('selectIds'),
          deleteSelectedRecords: createSpy('deleteSelectedRecords'),
          updateSelectedRecords: createSpy('updateSelectedRecords')
        }
      },
        {provide: RecordService, useValue: createSpyObj('RecordService', ['updatePages'])}
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

  describe('Record Selection', () => {
    it('should select single record', () => {
      const id = '1';
      component.onSelectRecords([id]);
      expect(inboxService.selectIds).toHaveBeenCalledWith([id]);
    });

    it('should add record in multi-select mode', () => {
      const id = '3';
      const selectedIds = ['1', '2'];

      inboxService.selectedIds$ = of(selectedIds);
      inboxService.isMultiSelect$ = of(true);
      component.onSelectRecords([id]);
      expect(inboxService.selectIds).toHaveBeenCalledWith([...selectedIds, id]);
    });

    it('should remove record in multi-select mode', () => {
      const id = '3';
      const selectedIds = ['1', '2', '3'];

      inboxService.selectedIds$ = of(selectedIds);
      inboxService.isMultiSelect$ = of(true);
      component.onSelectRecords([id]);
      expect(inboxService.selectIds).toHaveBeenCalledWith(['1', '2']);
    });
  });

  it('should upload pdfs on drop', () => {
    const files = [
      {},
      {},
      {}
    ];
    const event = {
      dataTransfer: {
        files
      }
    } as unknown as DragEvent;

    component.onDrop(event);
    expect(inboxService.upload).toHaveBeenCalledTimes(3);
  });
});
