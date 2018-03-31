import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {environment} from "../../environments/environment";

@Injectable()
export class TagService {
  constructor(private http: HttpClient) {
  }

  public get() {
    return this.http.get<string[]>(environment.api + "/tags");
  }

  public getPrimaryTags() {
    return this.http.get<string[]>(environment.api + "/categories");
  }

  public getByPatientId(id: string) {
    return this.http.get<string[]>(`${environment.api}/patients/${id}/tags`);
  }
}
