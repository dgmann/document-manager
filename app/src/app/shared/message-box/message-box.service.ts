import { Injectable } from "@angular/core";
import { MatDialog } from "@angular/material";
import { MessageBoxComponent } from "./message-box.component";
import { Observable } from "rxjs";

@Injectable({
  providedIn: "root"
})
export class MessageBoxService {
  constructor(private dialog: MatDialog) {
  }

  public open(title: string, text: string): Observable<boolean> {
    return this.dialog.open(MessageBoxComponent, {
      data: {
        title,
        text
      }
    }).afterClosed();
  }
}
