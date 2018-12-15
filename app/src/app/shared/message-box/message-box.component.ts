import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA } from "@angular/material";

@Component({
  selector: 'app-message-box',
  templateUrl: './message-box.component.html',
  styleUrls: ['./message-box.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MessageBoxComponent implements OnInit {
  title = "";
  text = "";

  constructor(@Inject(MAT_DIALOG_DATA) public data: MessageBoxData) {
    this.title = data.title;
    this.text = data.text;
  }

  ngOnInit() {
  }

}

export interface MessageBoxData {
  title: string;
  text: string
}
