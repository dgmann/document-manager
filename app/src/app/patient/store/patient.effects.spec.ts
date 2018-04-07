import {TestBed} from '@angular/core/testing';
import {provideMockActions} from '@ngrx/effects/testing';
import {Observable} from 'rxjs/Observable';

import {PatientEffects} from './patient.effects';

describe('ExternalApiService', () => {
  let actions$: Observable<any>;
  let effects: PatientEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PatientEffects,
        provideMockActions(() => actions$)
      ]
    });

    effects = TestBed.get(PatientEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
