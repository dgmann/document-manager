import { Injectable } from "@angular/core";
import { MatDialog } from "@angular/material";
import { DocumentEditDialogComponent } from "./document-edit-dialog.component";
import { Record } from "../../core/store";
import { filter, map } from "rxjs/operators";
import { Moment } from "moment";

@Injectable({
  providedIn: "root"
})
export class DocumentEditDialogService {

  constructor(private dialog: MatDialog) {
  }

  open(record: Record) {
    return this.dialog.open(DocumentEditDialogComponent, {
      disableClose: true,
      data: record,
      width: "635px"
    }).afterClosed()
      .pipe(
        filter(record => !!record),
        map(result => ({
          id: result.id, change: {
            patientId: result.patientId,
            date: result.date,
            tags: result.tags,
            category: result.category
          }
        } as EditResult))
      );
  }
}

export interface EditResult {
  id: string,
  change: {
    patientId: string,
    date: Moment,
    tags: string[],
    category: string
  }
}
