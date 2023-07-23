import { NgForOf, TitleCasePipe } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { Category, MatchType } from '@app/core/categories';

@Component({
  selector: 'app-categorydialog',
  templateUrl: './category-dialog.component.html',
  styleUrls: ['./category-dialog.component.scss'],
  standalone: true,
  imports: [MatDialogModule, MatFormFieldModule, MatInputModule, FormsModule, MatButtonModule, MatSelectModule, TitleCasePipe, NgForOf, MatCheckboxModule],
})
export class CategoryDialogComponent {
  typeOptions: string[] = [
    MatchType.Exact,
    MatchType.Regex,
  ];

  category: Category = {id: "", name: "", match: {type: MatchType.None, expression: ""}}

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: Category,
  ) {
    if (data != null) {
      this.category = data;
    }
  }

  protected readonly MatchType = MatchType;

  disableMatch(checked: boolean) {
    if (checked) {
      this.category.match.type = MatchType.None;
    } else {
      this.category.match.type = MatchType.Exact;
    }
  }
}
