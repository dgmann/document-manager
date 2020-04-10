import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {
  private appConfig;

  constructor(private http: HttpClient) {
  }

  public loadAppConfig() {
    return this.http
      .get('/assets/config.json')
      .toPromise()
      .then(data => {
        this.appConfig = data;
      });
  }

  public getApiUrl() {
    return `${this.getProto('http')}${this.appConfig.host}/api`;
  }

  public getNotificationWebsocketUrl() {
    return `${this.getProto('ws')}${this.appConfig.host}/api/notifications`;
  }

  private getProto(base: string) {
    let schema = base;
    if (this.appConfig.useSSL) {
      schema = base + 's';
    }
    return `${schema}://`;
  }
}
