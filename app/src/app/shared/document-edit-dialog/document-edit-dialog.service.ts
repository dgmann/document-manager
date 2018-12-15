import { Injectable } from "@angular/core";
import { MatDialog } from "@angular/material";
import { DocumentEditDialogComponent } from "./document-edit-dialog.component";
import { Record } from "../../core/store";
import { Observable } from "rxjs";
import { EditResult } from "./edit-result.model";
import { filter } from "rxjs/operators";

@Injectable({
  providedIn: "root"
})
export class DocumentEditDialogService {
  constructor(private dialog: MatDialog) {
  }

  open(record: Record): Observable<EditResult> {
    return this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record,
      width: "635px"
    }).afterClosed().pipe(
      filter(result => !!result)
    );
  }
}
