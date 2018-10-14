import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiRecordListComponent } from './multi-record-list.component';
import { MatCardModule, MatIconModule, MatTabsModule } from "@angular/material";
import { DocumentEditDialogService } from "../../shared";
import { of } from "rxjs";
import createSpyObj = jasmine.createSpyObj;

describe('MultiRecordListComponent', () => {
  let component: MultiRecordListComponent;
  let fixture: ComponentFixture<MultiRecordListComponent>;
  let editService;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatTabsModule,
        MatCardModule,
        MatIconModule
      ],
      declarations: [MultiRecordListComponent],
      providers: [
        {provide: DocumentEditDialogService, useValue: createSpyObj(["open"])}
      ]
    })
      .compileComponents();

  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiRecordListComponent);
    component = fixture.componentInstance;
    component.records = of([]);
    component.selectedCategory = of("1");
    component.categories = of({});
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
