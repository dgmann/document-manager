import { PdfLinkPipe } from './pdf-link.pipe';

describe('PdfLinkPipe', () => {
  let pipe: PdfLinkPipe;

  beforeEach(() => {
    const configService = jasmine.createSpyObj('ConfigService', {
      getApiUrl: '1'
    });
    pipe = new PdfLinkPipe('de-DE', configService);
  });

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('handles invalid values', () => {
    expect(pipe.transform(undefined)).toBe('');
  });
});
