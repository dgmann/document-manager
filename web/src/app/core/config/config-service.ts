import { DOCUMENT } from '@angular/common';
import { Inject, Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {
  private appConfig;

  constructor(private http: HttpClient, @Inject(DOCUMENT) private document: Document) {
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
    if (this.appConfig.host == null || this.appConfig.host === '') {
      return `${this.document.location.protocol}//${this.document.location.host}/api`;
    }
    return `${this.getProto('http')}${this.appConfig.host}/api`;
  }

  public getNotificationWebsocketUrl() {
    if (this.appConfig.host == null || this.appConfig.host === '') {
      const proto = this.document.location.protocol === 'https:' ? 'wss:' : 'ws:';
      return `${proto}//${this.document.location.host}/api/notifications`;
    }
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
