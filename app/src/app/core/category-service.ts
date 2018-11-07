import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "../../environments/environment";
import { map } from "rxjs/operators";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable({
  providedIn: "root"
})
export class CategoryService {
  public categories: Observable<Category[]>;
  public categoryMap: Observable<{ [id: string]: Category }>;

  constructor(private http: HttpClient) {
    this.categories = new BehaviorSubject<Category[]>([]);
    this.categoryMap = this.categories.pipe(
      map(cat => cat.reduce((accumulator, currentValue) => ({
        ...accumulator,
        [currentValue.id]: currentValue
      }), {} as { [id: string]: Category }))
    );
  }

  public load() {
    this.http.get<Category[]>(environment.api + "/categories").subscribe(categories => (this.categories as BehaviorSubject<Category[]>).next(categories));
  }

  public getByPatientId(id: string) {
    return this.http.get<Category[]>(`${environment.api}/patients/${id}/categories`);
  }
}

export interface Category {
  id: string;
  name: string;
}
