import {Injectable} from '@angular/core';
import {uniq} from 'lodash-es';

@Injectable({
  providedIn: 'root'
})
export class MultiSelectService<T> {
  private selection: T[] = [];

  public select(item: T, dataSource: T[], event: MouseEvent) {
    if (event.getModifierState('Control')) {
      const index = this.selection.indexOf(item);

      if (index >= 0) {
        this.selection.splice(index, 1);
      } else {
        this.selection = uniq([...this.selection, item]);
      }
    } else if (event.getModifierState('Shift')) {
      if (this.selection.length === 0) {
        this.selection = [item];
        return this.selection;
      }
      const firstCategory = this.selection[0];
      const selectFrom = dataSource.indexOf(firstCategory);
      const selectUntil = dataSource.indexOf(item);
      const indices = [selectFrom, selectUntil].sort();
      this.selection = dataSource.slice(indices[0], indices[1] + 1);
    } else {
      if (this.selection.length === 1 && this.selection[0] && this.selection[0] === item) {
        this.selection = [];
      } else {
        this.selection = [item];
      }
    }
    return this.selection;
  }
}
