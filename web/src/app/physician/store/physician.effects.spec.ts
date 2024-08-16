import {TestBed} from '@angular/core/testing';
import {provideMockActions} from '@ngrx/effects/testing';

import {PhysicianEffects} from './physician.effects';

describe('InboxService', () => {
  let effects: PhysicianEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PhysicianEffects,
        provideMockActions(() => null)
      ]
    });

    effects = TestBed.get(PhysicianEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
