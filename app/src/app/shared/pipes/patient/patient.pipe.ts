import {formatDate} from '@angular/common';
import {Inject, LOCALE_ID, Pipe, PipeTransform} from '@angular/core';
import {Patient} from '@app/patient';

@Pipe({
  name: 'patient'
})
export class PatientPipe implements PipeTransform {
  public static WITH_DATE_OF_BIRTH = 'withDateOfBirth';
  public static LONG = 'long';

  transform(patient: Patient, format?: string): unknown {
    if (!patient) {
      return '';
    }
    let res = `${patient.lastName}, ${patient.firstName}`;
    if (format === PatientPipe.LONG) {
      res = `${patient.id}: ${res}`;
    }
    if (format === PatientPipe.WITH_DATE_OF_BIRTH || format === PatientPipe.LONG) {
      res = `${res} (${formatDate(patient.birthDate, 'mediumDate', this.locale)})`;
    }
    return res;
  }

  constructor(@Inject(LOCALE_ID) public locale: string){}
}
