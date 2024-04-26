import {CategoryPipe} from './category.pipe';
import {of} from 'rxjs';
import {CategoryService} from '@app/core/categories';

describe('CategoryPipePipe', () => {
  it('create an instance', () => {
    const result = of({});
    const pipe = new CategoryPipe({
      categoryMap: result
    } as CategoryService);
    expect(pipe).toBeTruthy();
  });
});
