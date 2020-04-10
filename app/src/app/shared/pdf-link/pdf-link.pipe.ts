import {Inject, LOCALE_ID, Pipe, PipeTransform} from '@angular/core';
import {DatePipe} from '@angular/common';
import {ConfigService} from '@app/core/config';

@Pipe({
  name: 'pdfLink'
})
export class PdfLinkPipe implements PipeTransform {

  constructor(@Inject(LOCALE_ID) public locale: string, private config: ConfigService) {
  }

  transform(ids: string | string[], title?: string): any {
    if (typeof ids === 'string') {
      ids = [ids];
    }
    if (!title) {
      title = `Exportiert um: ${new DatePipe(this.locale).transform(new Date(), 'short')}`;
    }
    const url = new URL(`${this.config.getApiUrl()}/export`);
    ids.forEach(id => url.searchParams.append('id', id));
    url.searchParams.append('title', title);
    return url.href;
  }

}
