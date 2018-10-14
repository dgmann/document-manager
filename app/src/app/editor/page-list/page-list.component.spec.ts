import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PageListComponent } from './page-list.component';
import { MatButtonModule, MatCardModule, MatIconModule } from "@angular/material";
import { DndModule } from "ng2-dnd";

describe('PageListComponent', () => {
  let component: PageListComponent;
  let fixture: ComponentFixture<PageListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatCardModule,
        MatButtonModule,
        MatIconModule,
        DndModule.forRoot(),
      ],
      declarations: [PageListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PageListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
