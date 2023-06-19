import {NgModule} from '@angular/core';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {EffectsModule} from '@ngrx/effects';
import {StoreModule} from '@ngrx/store';
import {SharedModule} from '../shared';
import {InboxComponent} from './inbox.component';
import {InboxRouterModule} from './inbox.routes';
import {metaReducers, reducers} from './reducers';
import {InboxEffects} from './store/inbox.effects';
import {MatBottomSheetModule} from '@angular/material/bottom-sheet';
import {ActionBarComponent} from '@app/inbox/action-bar/action-bar.component';
import {MatIconModule} from '@angular/material/icon';
import {ActionBarService} from '@app/inbox/action-bar/action-bar.service';
import {MatMenuModule} from '@angular/material/menu';
import {MatListModule} from '@angular/material/list';
import {MatButtonModule} from '@angular/material/button';

@NgModule({
    imports: [
        StoreModule.forFeature('inbox', reducers, { metaReducers }),
        EffectsModule.forFeature([InboxEffects]),
        InboxRouterModule,
        MatSlideToggleModule,
        MatProgressSpinnerModule,
        MatBottomSheetModule,
        SharedModule,
        MatIconModule,
        MatMenuModule,
        MatListModule,
        MatButtonModule
    ],
    declarations: [
        InboxComponent,
        ActionBarComponent
    ],
    exports: [
        InboxComponent
    ],
    providers: [
        ActionBarService
    ]
})
export class InboxModule {
}
