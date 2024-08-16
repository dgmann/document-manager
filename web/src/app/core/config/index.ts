import { Provider } from '@angular/core';
import { ConfigService } from '@app/core/config/config-service';

export * from './config-service';

export function provideMockConfigService(): Provider {
  const service = {
    getApiUrl: () => 'http://test.com',
    loadAppConfig: () => null,
    getNotificationWebsocketUrl: () => 'http://test.com/ws'
  };
  return { provide: ConfigService, useValue: service };
}
