import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Observable } from "rxjs";
import { Category, CategoryService } from "../../shared/category-service";
import { Record, RecordService } from "../../store";
import { DocumentEditDialogComponent } from "../../shared/document-edit-dialog/document-edit-dialog.component";
import { MatDialog } from "@angular/material";

@Component({
  selector: 'app-multi-record-list',
  templateUrl: './multi-record-list.component.html',
  styleUrls: ['./multi-record-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MultiRecordListComponent implements OnInit {
  @Input('records') records: Observable<Record[]>;
  @Output('clickRecord') clickRecord = new EventEmitter<string>();
  categories: Observable<{ [id: string]: Category }>;

  constructor(private categoryService: CategoryService,
              private recordService: RecordService,
              private dialog: MatDialog) {
  }

  ngOnInit() {
    this.categories = this.categoryService.getAsMap()
  }

  onRecordClicked(id: string) {
    this.clickRecord.emit(id);
  }

  edit(record: Record) {
    this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record,
      width: "635px"
    }).afterClosed().subscribe((result: Record) => {
      if (!result) {
        return;
      }
      this.recordService.update(result.id, {
        patientId: result.patientId,
        date: result.date,
        tags: result.tags,
        categoryId: result.categoryId
      });
    });
  }

}
