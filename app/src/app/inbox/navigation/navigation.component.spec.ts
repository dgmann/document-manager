import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {NavigationComponent} from './navigation.component';
import {NgDragDropModule} from 'ng-drag-drop';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {InboxService} from '../inbox.service';
import {of} from 'rxjs';
import createSpy = jasmine.createSpy;

describe('Inbox NavigationComponent', () => {
  let component: NavigationComponent;
  let fixture: ComponentFixture<NavigationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        NgDragDropModule,
        MatSlideToggleModule,
      ],
      declarations: [NavigationComponent],
      providers: [{
        provide: InboxService, useValue: {
          isMultiSelect$: of(false),
          setMultiselect: createSpy('setMultiselect')
        }
      }]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NavigationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
