import {Injectable} from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import {of} from 'rxjs';
import {catchError, filter, switchMap, take, tap} from 'rxjs/operators';
import {RecordService} from '@app/core/records';

@Injectable()
export class EditorGuard  {
  constructor(private recordService: RecordService) {
  }

  getFromStoreOrAPI(id: string) {
    return this.recordService.find(id)
      .pipe(
        tap(record => {
          if (!record) {
            this.recordService.load({id});
          }
        }),
        filter(record => !!record),
        take(1));
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    return this.getFromStoreOrAPI(route.params.id).pipe(
      switchMap(() => of(true)),
      catchError(() => of(false))
    );
  }
}
