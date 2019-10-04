import {MatBottomSheet, MatBottomSheetConfig} from '@angular/material/bottom-sheet';
import {ActionBarComponent} from '@app/inbox/action-bar/action-bar.component';
import {Injectable} from '@angular/core';

@Injectable()
export class ActionBarService {
  public isOpen: boolean;

  constructor(private bottomSheet: MatBottomSheet) {
  }

  public open(config?: MatBottomSheetConfig) {
    this.bottomSheet.open(ActionBarComponent, {
      hasBackdrop: false,
      ...config
    });
    this.isOpen = true;
  }

  public dismiss() {
    this.bottomSheet.dismiss();
    this.isOpen = false;
  }
}
