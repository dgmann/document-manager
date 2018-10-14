import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditorComponent } from './editor.component';
import { PageListComponent } from "./page-list/page-list.component";
import { MatButtonModule, MatCardModule, MatIconModule } from "@angular/material";
import { DndModule } from "ng2-dnd";
import { RecordService } from "../core/store";
import { RouterTestingModule } from "@angular/router/testing";
import createSpyObj = jasmine.createSpyObj;

describe('EditorComponent', () => {
  let component: EditorComponent;
  let fixture: ComponentFixture<EditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatCardModule,
        MatButtonModule,
        MatIconModule,
        DndModule.forRoot(),
        RouterTestingModule
      ],
      declarations: [
        EditorComponent,
        PageListComponent
      ],
      providers: [
        {provide: RecordService, useValue: createSpyObj(['find', 'updatePages', 'getInvalidIds'])}
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
