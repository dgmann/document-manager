import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'ids'
})
export class IdsPipe implements PipeTransform {

  transform(values: {id: string}[]): string[] {
    if (!values) {
      return [];
    }
    return values.map(v => v.id);
  }

}
