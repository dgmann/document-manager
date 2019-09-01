import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {environment} from '@env/environment';
import {BehaviorSubject, Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TagService {
  public tags: Observable<string[]>;

  constructor(private http: HttpClient) {
    this.tags = new BehaviorSubject<string[]>([]);
  }

  public load() {
    this.http.get<string[]>(environment.api + '/tags').subscribe(tags => (this.tags as BehaviorSubject<string[]>).next(tags));
  }

  public getByPatientId(id: string) {
    return this.http.get<string[]>(`${environment.api}/patients/${id}/tags`);
  }
}
