import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CaseComponent } from './case/case.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { MakeCaseComponent } from './make-case/make-case.component';
import { MapTestComponent } from './map-test/map-test.component';

const routes: Routes = [
  { path: 'home', component:HomeComponent },
  { path: 'login', component:LoginComponent},
  { path: 'dashboard', component:DashboardComponent},
  { path: 'case', component:CaseComponent},
  { path: 'makeCase', component:MakeCaseComponent},
  { path: 'testMap', component:MapTestComponent},
  { path: '**', component:LoginComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
