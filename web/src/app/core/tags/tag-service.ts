import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {environment} from '@env/environment';
import {BehaviorSubject, Observable} from 'rxjs';
import {ConfigService} from '@app/core/config/config-service';

@Injectable({
  providedIn: 'root'
})
export class TagService {
  public tags: Observable<string[]>;

  constructor(private http: HttpClient, private config: ConfigService) {
    this.tags = new BehaviorSubject<string[]>([]);
  }

  public load() {
    this.http.get<string[]>(this.config.getApiUrl() + '/tags').subscribe(tags => (this.tags as BehaviorSubject<string[]>).next(tags));
  }

  public getByPatientId(id: string) {
    return this.http.get<string[]>(`${this.config.getApiUrl()}/patients/${id}/tags`);
  }
}
