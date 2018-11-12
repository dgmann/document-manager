import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditorComponent } from './editor.component';
import { PageListComponent } from "./page-list/page-list.component";
import { MatButtonModule, MatCardModule, MatIconModule } from "@angular/material";
import { RecordService } from "../core/store";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";
import createSpy = jasmine.createSpy;

describe('EditorComponent', () => {
  let component: EditorComponent;
  let fixture: ComponentFixture<EditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatCardModule,
        MatButtonModule,
        MatIconModule,
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
