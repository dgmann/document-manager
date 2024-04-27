import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'contains'
})
export class ContainsPipe implements PipeTransform {

  transform<T>(value: T[], arg: T): boolean {
    return value.includes(arg);
  }

}
