import {TestBed} from '@angular/core/testing';
import {provideMockActions} from '@ngrx/effects/testing';

import {InboxEffects} from './inbox.effects';

describe('InboxEffects', () => {
  let effects: InboxEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        InboxEffects,
        provideMockActions(() => null)
      ]
    });

    effects = TestBed.get(InboxEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
