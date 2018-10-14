import { TestBed } from '@angular/core/testing';
import { provideMockActions } from '@ngrx/effects/testing';
import { Observable } from 'rxjs/Observable';

import { RecordEffects } from './record.effects';
import { HttpClientTestingModule } from "@angular/common/http/testing";

describe('RecordService', () => {
  let actions$: Observable<any>;
  let effects: RecordEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        HttpClientTestingModule
      ],
      providers: [
        RecordEffects,
        provideMockActions(() => actions$)
      ]
    });

    effects = TestBed.get(RecordEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
