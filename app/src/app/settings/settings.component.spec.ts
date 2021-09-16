import { HttpClientTestingModule } from '@angular/common/http/testing';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {MatInputModule} from '@angular/material/input';
import {NoopAnimationsModule} from '@angular/platform-browser/animations';
import {ConfigService} from '@app/core/config';
import {RecordService} from '@app/core/records';
import {PatientService} from '@app/patient/patient.service';
import {of} from 'rxjs';

import { SettingsComponent } from './settings.component';
import {MatCardModule} from '@angular/material/card';
import {MatListModule} from '@angular/material/list';
import {FormsModule} from '@angular/forms';
import {MatFormFieldModule} from '@angular/material/form-field';

describe('SettingsComponent', () => {
  let component: SettingsComponent;
  let fixture: ComponentFixture<SettingsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        HttpClientTestingModule,
        MatCardModule,
        MatListModule,
        MatFormFieldModule,
        MatInputModule,
        FormsModule,
        NoopAnimationsModule],
      providers: [
        { provide: PatientService, useValue: {categories: of(), load: () => {} }},
        { provide: ConfigService, useValue: {getApiUrl: () => 'http://test.com' }},
      ],
      declarations: [ SettingsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
