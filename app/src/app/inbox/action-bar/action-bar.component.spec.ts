import {HttpClientTestingModule} from '@angular/common/http/testing';
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {InboxService} from '@app/inbox/inbox.service';
import {SharedModule} from '@app/shared';
import {of} from 'rxjs';

import {ActionBarComponent} from './action-bar.component';
import {MatIconModule} from '@angular/material/icon';

describe('EventSnackbarComponent', () => {
  let component: ActionBarComponent;
  let fixture: ComponentFixture<ActionBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatIconModule,
        HttpClientTestingModule,
        MatSnackBarModule,
        SharedModule
      ],
      declarations: [ ActionBarComponent ],
      providers: [{ provide: InboxService, useValue: { selectedIds$: of() } }]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
