import {Pipe, PipeTransform} from '@angular/core';
import {CategoryService} from '@app/core';
import {map} from 'rxjs/operators';
import {Observable} from 'rxjs';

@Pipe({
  name: 'category',
  pure: false
})
export class CategoryPipe implements PipeTransform {

  constructor(private categoryService: CategoryService) {
  }

  transform(value: string, ...args: any[]): Observable<string> {
    return this.categoryService.categoryMap.pipe(map(categories => categories[value].name));
  }

}
