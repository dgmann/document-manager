import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatMenuModule } from '@angular/material/menu';

import {EditorComponent} from './editor.component';
import {PageListComponent} from './page-list/page-list.component';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {RecordService} from '../core/records';
import {RouterTestingModule} from '@angular/router/testing';
import {of} from 'rxjs';
import createSpy = jasmine.createSpy;

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
        RouterTestingModule
      ],
      declarations: [
        EditorComponent,
        PageListComponent
      ],
      providers: [
        {
          provide: RecordService, useValue: {
            find: createSpy().and.returnValue(of({pages: []})),
            updatePages: createSpy(),
            getInvalidIds: createSpy()
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
