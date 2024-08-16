import {Patient} from '@app/patient';
import { PatientPipe } from './patient.pipe';
import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';
import localeDeExtra from '@angular/common/locales/extra/de';


describe('PatientPipe', () => {
  beforeAll(() => {
    registerLocaleData(localeDe, 'de-DE', localeDeExtra);
  })

  const patient: Patient = {
    birthDate: new Date('01.01.1900'),
    firstName: 'Test',
    lastName: 'Person',
    id: '1'
  };

  it('create an instance', () => {
    const pipe = new PatientPipe('de-DE');
    expect(pipe).toBeTruthy();
  });

  it('formats with birthDate', () => {
    const pipe = new PatientPipe('de-DE');
    expect(pipe.transform(patient, PatientPipe.WITH_DATE_OF_BIRTH)).toEqual('Person, Test (01.01.1900)');
  });

  it('formats without birthDate', () => {
    const pipe = new PatientPipe('de-DE');
    expect(pipe.transform(patient)).toEqual('Person, Test');
  });
});
