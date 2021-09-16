import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {ConfigService} from '@app/core/config';
import {RecordService} from '@app/core/records';
import {HistoryService} from '@app/history/history-service';
import { HistoryModule } from '@app/history/history.module';
import {of} from 'rxjs';

import {HistoryComponent} from './history.component';

describe('HistoryComponent', () => {
  let component: HistoryComponent;
  let fixture: ComponentFixture<HistoryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [ {provide: HistoryService, useValue: {
        next: () => {},
        records$: of(),
        selectedId$: of(),
        selectedRecord$: of()
      }}],
      declarations: [HistoryComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HistoryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
