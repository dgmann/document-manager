import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "../../environments/environment";
import { map } from "rxjs/operators";

@Injectable()
export class CategoryService {
  constructor(private http: HttpClient) {
  }

  public get() {
    return this.http.get<Category[]>(environment.api + "/categories");
  }

  public getAsMap() {
    return this.get().pipe(
      map(cat => cat.reduce((accumulator, currentValue) => ({
        ...accumulator,
        [currentValue.id]: currentValue
      }), {} as { [id: string]: Category }))
    );
  }

  public getByPatientId(id: string) {
    return this.http.get<Category[]>(`${environment.api}/patients/${id}/categories`);
  }
}

export interface Category {
  id: string;
  name: string;
}
