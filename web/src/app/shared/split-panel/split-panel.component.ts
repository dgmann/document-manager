import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-split-panel',
  templateUrl: './split-panel.component.html',
  styleUrls: ['./split-panel.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SplitPanelComponent implements OnInit {

  constructor() {
  }

  ngOnInit() {
  }

}
