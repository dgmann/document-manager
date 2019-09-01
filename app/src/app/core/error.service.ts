import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ErrorService {
  public errors$: Observable<string>;

  constructor() {
  }
}
