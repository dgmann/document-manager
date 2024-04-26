import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CategoryDialogComponent } from './category-dialog.component';

describe('CategorydialogComponent', () => {
  let component: CategoryDialogComponent;
  let fixture: ComponentFixture<CategoryDialogComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CategoryDialogComponent]
    });
    fixture = TestBed.createComponent(CategoryDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
