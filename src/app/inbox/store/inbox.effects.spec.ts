import {TestBed} from '@angular/core/testing';
import {provideMockActions} from '@ngrx/effects/testing';
import {Observable} from 'rxjs/Observable';

import {InboxEffects} from './inbox.effects';

describe('InboxService', () => {
  let actions$: Observable<any>;
  let effects: InboxEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        InboxEffects,
        provideMockActions(() => actions$)
      ]
    });

    effects = TestBed.get(InboxEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
