import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CategoryDialogComponent } from './category-dialog.component';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Category, MatchType } from '@app/core/categories';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('CategorydialogComponent', () => {
  let component: CategoryDialogComponent;
  let fixture: ComponentFixture<CategoryDialogComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [CategoryDialogComponent, NoopAnimationsModule],
      providers: [
        { provide: MAT_DIALOG_DATA, useValue: {
          id: "test", 
          match: {
            expression: "", 
            type: MatchType.All
          }, 
          name: "Test" } as Category 
        }
      ]
    });
    fixture = TestBed.createComponent(CategoryDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
