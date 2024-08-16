import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatMenuModule } from '@angular/material/menu';

import {EditorComponent} from './editor.component';
import {PageListComponent} from './page-list/page-list.component';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {RecordService} from '../core/records';
import {of} from 'rxjs';
import { RouterModule } from '@angular/router';
import { DragDropModule } from '@angular/cdk/drag-drop';

describe('EditorComponent', () => {
  let component: EditorComponent;
  let fixture: ComponentFixture<EditorComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        MatCardModule,
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
        RouterModule.forRoot([]),
        DragDropModule
      ],
      declarations: [
        EditorComponent,
        PageListComponent
      ],
      providers: [
        {
          provide: RecordService, useValue: {
            find: jest.fn().mockReturnValue(of({pages: []})),
            updatePages: jest.fn(),
            getInvalidIds: jest.fn()
          }
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
