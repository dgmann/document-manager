import {Inject, LOCALE_ID, Pipe, PipeTransform} from '@angular/core';
import {environment} from '@env/environment';
import {DatePipe} from '@angular/common';

@Pipe({
  name: 'pdfLink'
})
export class PdfLinkPipe implements PipeTransform {

  constructor(@Inject(LOCALE_ID) public locale: string) {}

  transform(ids: string | string[], title?: string): any {
    if (typeof ids === 'string') {
      ids = [ids];
    }
    if (!title) {
      title = `Exportiert um: ${new DatePipe(this.locale).transform(new Date(), 'short')}`;
    }
    const url = new URL(`${environment.api}/export`);
    ids.forEach(id => url.searchParams.append('id', id));
    url.searchParams.append('title', title);
    return url.href;
  }

}
