import {HttpClient} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import {BehaviorSubject, Observable} from 'rxjs';
import {ConfigService} from '@app/core/config/config-service';
import {Category} from './category';

@Injectable({
  providedIn: 'root'
})
export class CategoryService {
  public categories: Observable<Category[]>;
  public categoryMap: Observable<{ [id: string]: Category }>;

  constructor(private http: HttpClient, private config: ConfigService) {
    this.categories = new BehaviorSubject<Category[]>([]);
    this.categoryMap = this.categories.pipe(
      map(cat => cat.reduce((accumulator, currentValue) => ({
        ...accumulator,
        [currentValue.id]: currentValue
      }), {} as { [id: string]: Category }))
    );
  }

  public load() {
    this.http.get<Category[]>(this.config.getApiUrl() + '/categories')
      .subscribe(categories => (this.categories as BehaviorSubject<Category[]>).next(categories));
  }

  public getByPatientId(id: string) {
    return this.http.get<Category[]>(`${this.config.getApiUrl()}/patients/${id}/categories`);
  }

  public add(category: Category) {
    return this.http.post(`${this.config.getApiUrl()}/categories`, category);
  }

  public update(category: Category) {
    return this.http.put(`${this.config.getApiUrl()}/categories`, category);
  }
}
