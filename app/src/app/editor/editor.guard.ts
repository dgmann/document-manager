import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot } from "@angular/router";
import { RecordService } from "../store";
import { catchError, filter, switchMap, take, tap } from "rxjs/operators";
import { of } from "rxjs/observable/of";

@Injectable()
export class EditorGuard implements CanActivate {
  constructor(private recordService: RecordService) {
  }

  getFromStoreOrAPI(id: string) {
    return this.recordService.find(id)
      .pipe(
        tap(record => {
          if (!record) {
            this.recordService.load({id: id})
          }
        }),
        filter(record => !!record),
        take(1));
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    return this.getFromStoreOrAPI(route.params['id']).pipe(
      switchMap(() => of(true)),
      catchError(() => of(false))
    );
  }
}