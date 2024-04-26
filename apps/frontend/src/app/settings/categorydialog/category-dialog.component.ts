import { NgForOf } from '@angular/common';
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
  imports: [MatDialogModule, MatFormFieldModule, MatInputModule, FormsModule, MatButtonModule, MatSelectModule, NgForOf, MatCheckboxModule],
})
export class CategoryDialogComponent {
  typeOptions: {type: MatchType, display: string}[] = [
    {type: MatchType.All, display: "Alle Wörter"},
    {type: MatchType.Any, display: "Irgendein Wort"},
    {type: MatchType.Exact, display: "Exakt"},
    {type: MatchType.Regex, display: "Regex"},
  ];
  typeDescriptions = {
    [MatchType.All]: "Dokument enthält alle Wörter (getrennt durch Leerzeichen)",
    [MatchType.Any]: "Dokument enthält irgendeines der Wörter (getrennt durch Leerzeichen)",
    [MatchType.Exact]: "Dokument enthält exakt die Zeichenfolge",
    [MatchType.Regex]: "Dokument passt zu dem Ausdruck",
  }

  category: Category = {id: "", name: "", match: {type: MatchType.None, expression: ""}}

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: Category,
  ) {
    if (data != null) {
      this.category = data;
    }
  }

  protected readonly MatchType = MatchType;

  automatchingChecked(checked: boolean) {
    if (checked) {
      this.category.match.type = MatchType.Exact;
    } else {
      this.category.match.type = MatchType.None;
    }
  }
}
