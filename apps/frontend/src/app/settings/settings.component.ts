import {Component, OnInit} from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {Category, CategoryService} from '@app/core/categories';
import { CategoryDialogComponent } from '@app/settings/categorydialog/category-dialog.component';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  categories$: Observable<Category[]>;

  constructor(private categoryService: CategoryService, public dialog: MatDialog) {
    this.categories$ = this.categoryService.categories;
  }

  ngOnInit() {
    this.categoryService.load();
  }

  edit(category: Category) {
    const dialogRef = this.dialog.open(CategoryDialogComponent, {
      data: category,
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.categoryService.update(result).subscribe(() => this.categoryService.load());
      }
    });
  }

  add() {
    const dialogRef = this.dialog.open(CategoryDialogComponent);

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.categoryService.add(result).subscribe(() => this.categoryService.load());
      }
    });
  }
}
