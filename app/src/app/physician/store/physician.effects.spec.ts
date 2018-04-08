import {TestBed} from '@angular/core/testing';
import {provideMockActions} from '@ngrx/effects/testing';
import {Observable} from 'rxjs/Observable';

import {PhysicianEffects} from './physician.effects';

describe('InboxService', () => {
  let actions$: Observable<any>;
  let effects: PhysicianEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PhysicianEffects,
        provideMockActions(() => actions$)
      ]
    });

    effects = TestBed.get(PhysicianEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
