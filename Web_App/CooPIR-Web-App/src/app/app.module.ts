import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { MatIconModule } from '@angular/material/icon';
import { LoginComponent } from './login/login.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { MatToolbarModule } from '@angular/material/toolbar';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatMenuModule } from '@angular/material/menu';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatList, MatListModule } from '@angular/material/list';
import { MatSidenavModule } from '@angular/material/sidenav';
import { CookieService } from 'ngx-cookie-service';
import { httpInterceptProviders } from './Http-interceptors';
import { CaseComponent } from './case/case.component';
import { NgxDocViewerModule } from 'ngx-doc-viewer';
import * as FileSaver from 'file-saver';
import { MakeCaseComponent } from './make-case/make-case.component';
import { MapTestComponent } from './map-test/map-test.component';
import { JobPageComponent } from './job-page/job-page.component';
import { AgentPageComponent } from './agent-page/agent-page.component';
import { MapPageComponent } from './map-page/map-page.component';
import { Ng2SearchPipeModule } from 'ng2-search-filter';
import { FormsModule } from '@angular/forms';


@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LoginComponent,
    DashboardComponent,
    CaseComponent,
    MakeCaseComponent,
    MapTestComponent,
    JobPageComponent,
    AgentPageComponent,
    MapPageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatIconModule,
    HttpClientModule,
    MatToolbarModule,
    FlexLayoutModule,
    MatMenuModule,
    MatButtonModule,
    MatDividerModule,
    MatListModule,
    NgxDocViewerModule,
    Ng2SearchPipeModule,
    FormsModule
  ],
  providers: [
    CookieService,
    httpInterceptProviders
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
