import {Component, OnInit} from '@angular/core';
import {Category, CategoryService} from '@app/core/categories';
import {Observable} from 'rxjs';
import {NgForm} from '@angular/forms';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  categories$: Observable<Category[]>;

  constructor(private categoryService: CategoryService) {
    this.categories$ = this.categoryService.categories;
  }

  ngOnInit() {
    this.categoryService.load();
  }

  public addCategory(form: NgForm) {
    this.categoryService.add(form.value.id, form.value.name).subscribe(() => this.categoryService.load());
    form.resetForm();
  }

}
