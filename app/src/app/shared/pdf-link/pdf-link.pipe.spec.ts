import { PdfLinkPipe } from './pdf-link.pipe';

describe('PdfLinkPipe', () => {
  it('create an instance', () => {
    const configService = jasmine.createSpyObj('ConfigService', ['getApiUrl']);
    configService.and.returnValue('');
    const pipe = new PdfLinkPipe('de-DE', configService);
    expect(pipe).toBeTruthy();
  });
});
