import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {environment} from "../../environments/environment";

@Injectable()
export class CategoryService {
  constructor(private http: HttpClient) {
  }

  public get() {
    return this.http.get<Category[]>(environment.api + "/categories");
  }

  public getByPatientId(id: string) {
    return this.http.get<Category[]>(`${environment.api}/patients/${id}/categories`);
  }
}

export interface Category {
  id: string;
  name: string;
}
